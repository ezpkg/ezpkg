package ggen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"

	"ezpkg.io/errorz"
)

func (ng *engine) start(cfg Config, patterns ...string) (_err error) {
	{
		for _, plugin := range cfg.Plugins {
			if err := ng.registerPlugin(plugin); err != nil {
				return err
			}
		}
	}
	{
		if len(patterns) == 0 {
			return errorz.Wrapf(nil, "no patterns")
		}
		if len(ng.plugins) == 0 {
			return errorz.Wrapf(nil, "no registered plugins")
		}
		if err := ng.validateConfig(&cfg); err != nil {
			return err
		}
		ng.xcfg = cfg
	}
	buildFlags := getBuildFlags(cfg.BuildTags)
	{
		mode := packages.NeedName | packages.NeedImports | packages.NeedDeps |
			packages.NeedFiles | packages.NeedCompiledGoFiles
		ng.pkgcfg = packages.Config{
			Mode:       mode,
			BuildFlags: buildFlags,
		}
		pkgs, err := packages.Load(&ng.pkgcfg, patterns...)
		if err != nil {
			return errorz.Wrapf(err, "can not load package: %v", err)
		}

		// populate cleanedFileNames
		cleanedFileNames := make(map[string]bool)
		for _, pl := range ng.enabledPlugins {
			input := GenerateFileNameInput{PluginName: pl.name}
			filename := ng.genFilename(input)
			cleanedFileNames[filename] = true
		}

		// list available packages
		availablePkgs := make([]*packages.Package, 0, len(pkgs))
		for _, pkg := range pkgs {
			pkgDir := getPackageDir(pkg)
			if pkgDir == "" {
				ng.logger.Info("no Go files found in package", "pkg", pkg)
				continue
			}
			availablePkgs = append(availablePkgs, pkg)
			if err = cleanDir(cleanedFileNames, pkgDir); err != nil {
				return err
			}
		}
		ng.cleanedFileNames = cleanedFileNames
		if cfg.CleanOnly {
			return nil
		}

		// populate collectedPackages, includes, srcMap
		if err = ng.collectPackages(availablePkgs); err != nil {
			return err
		}

		if ng.logger.Enabled(DebugLevel) {
			for _, pkg := range ng.collectedPackages {
				ng.logger.Debug("collected package", "pkg", pkg.PkgPath)
			}
		}
	}
	{
		sortedIncludedPackages := make([]includedPackage, 0, len(ng.includedPackages))
		for pkgPath, includedFlags := range ng.includedPackages {
			if slices.Index(includedFlags, true) >= 0 {
				sortedIncludedPackages = append(sortedIncludedPackages, includedPackage{pkgPath, includedFlags})
			}
		}
		sort.Slice(sortedIncludedPackages, func(i, j int) bool {
			return sortedIncludedPackages[i].PkgPath < sortedIncludedPackages[j].PkgPath
		})
		ng.sortedIncludedPackages = sortedIncludedPackages

		if ng.logger.Enabled(DebugLevel) {
			for _, pkg := range sortedIncludedPackages {
				ng.logger.Debug("included package", "pkg", pkg.PkgPath)
			}
		}
	}
	{
		pkgPatterns := make([]string, 0, len(ng.includedPatterns)+len(ng.sortedIncludedPackages))
		pkgPatterns = append(pkgPatterns, ng.includedPatterns...)
		for _, pkg := range ng.sortedIncludedPackages {
			pkgPatterns = append(pkgPatterns, pkg.PkgPath)
		}
		if ng.logger.Enabled(DebugLevel) {
			ng.logger.Debug("load all syntax from:")
			for _, p := range pkgPatterns {
				ng.logger.Debug("  " + p)
			}
		}
		if len(pkgPatterns) == 0 {
			fmt.Println("no packages for generating")
			return nil
		}
		pkgPatterns = append(pkgPatterns, builtinPath) // load builtin types

		ng.pkgcfg = packages.Config{
			Mode:       packages.LoadAllSyntax,
			BuildFlags: buildFlags,
			Overlay:    ng.srcMap,
		}
		pkgs, err := packages.Load(&ng.pkgcfg, pkgPatterns...)
		if err != nil {
			return errorz.Wrapf(err, "can not load package: %v", err)
		}

		// populate xinfo
		ng.xinfo = newExtendedInfo(ng.pkgcfg.Fset)
		packages.Visit(pkgs,
			func(pkg *packages.Package) bool {
				if cfg.Namespace != "" && !strings.HasPrefix(pkg.PkgPath, cfg.Namespace) {
					return true
				}
				if err2 := ng.xinfo.AddPackage(pkg); err2 != nil {
					_err = err2
					return false
				}
				return true
			}, nil)
		if _err != nil {
			return _err
		}

		// populate pkgMap
		packages.Visit(pkgs,
			func(pkg *packages.Package) bool {
				ng.pkgMap[pkg.PkgPath] = pkg
				ng.dir2pkg[GetPkgDir(pkg)] = pkg
				return true
			}, nil)

		// populate builtin types
		ng.builtinTypes = parseBuiltinTypes(ng.pkgMap[builtinPath])
		delete(ng.pkgMap, builtinPath)
	}
	{
		// populate generatedFiles
		for _, pl := range ng.enabledPlugins {
			wrapNg := &wrapEngine{
				embededLogger: embededLogger{ng.logger.With("plugin", pl.name)},
				engine:        ng,
				plugin:        pl,
			}
			if err := pl.plugin.Generate(wrapNg); err != nil {
				return errorz.Wrapf(err, "%v: %v", pl.name, err)
			}
			for _, gpkg := range wrapNg.pkgs {
				prt := gpkg.printer
				if prt != nil && prt.buf.Len() != 0 {
					// close the printer for writing to file, but only if there
					// are any bytes written
					if err := prt.Close(); err != nil {
						return err
					}
				}
			}
		}
	}
	{
		sort.Strings(ng.generatedFiles)
		fmt.Println("Generated files:")
		pwd, err := os.Getwd()
		must(err)
		for _, filename := range ng.generatedFiles {
			filename, err = filepath.Rel(pwd, filename)
			must(err)
			fmt.Printf("\t./%v\n", filename)
		}
		if err = ng.execGoimport(ng.generatedFiles); err != nil {
			return err
		}
	}
	return nil
}

func (ng *engine) collectPackages(pkgs []*packages.Package) error {
	collectedPackages, fileContents, err := collectPackages(ng.logger, pkgs, ng.cleanedFileNames)
	if err != nil {
		return err
	}
	sort.Slice(collectedPackages, func(i, j int) bool {
		return collectedPackages[i].PkgPath < collectedPackages[j].PkgPath
	})
	pkgMap := map[string][]bool{}
	for _, pl := range ng.enabledPlugins {
		filterNg := &filterEngine{
			embededLogger: embededLogger{ng.logger.With("plugin", pl.name)},
			ng:            ng,
			plugin:        pl,
			pkgs:          collectedPackages,
			pkgMap:        pkgMap,
			patterns:      &ng.includedPatterns,
		}
		if err = pl.plugin.Filter(filterNg); err != nil {
			return errorz.Wrapf(err, "plugin %v: %v", pl.name, err)
		}
	}
	ng.collectedPackages = collectedPackages
	ng.includedPackages = pkgMap
	ng.mapPkgDirectives = make(map[string][]Directive)
	for _, pkg := range collectedPackages {
		ng.mapPkgDirectives[pkg.PkgPath] = pkg.Directives
	}

	srcMap := make(map[string][]byte)
	for _, content := range fileContents {
		srcMap[content.Path] = content.Body
	}
	ng.srcMap = srcMap
	return nil
}

func getBuildFlags(buildTags []string) []string {
	var buildFlags = "-tags ggen"
	if len(buildTags) > 0 {
		buildFlags += "," + strings.Join(buildTags, ",")
	}
	return strings.Split(buildFlags, " ")
}

type fileContent struct {
	Path string
	Body []byte
}

func collectPackages(logger Logger, pkgs []*packages.Package, cleanedFileNames map[string]bool) (collectedPackages []filteringPackage, files []fileContent, _err error) {

	var wg0, wg sync.WaitGroup
	wg0.Add(2)

	fileCh := make(chan fileContent, 4)
	go func() {
		defer wg0.Done()
		for file := range fileCh {
			files = append(files, file)
		}
	}()

	// collect errors
	errCh := make(chan error, 4)
	var errs []error
	go func() {
		defer wg0.Done()
		for err := range errCh {
			errs = append(errs, err)
		}
	}()

	limit := make(chan struct{}, 16) // limit concurrency
	collectedPackages = make([]filteringPackage, len(pkgs))
	for i := range pkgs {
		limit <- struct{}{} // limit
		wg.Add(1)
		go func(i int, pkg *packages.Package) {
			defer func() { wg.Done(); <-limit }() // release limit
			directives, inlineDirectives, err := parseDirectivesFromPackage(logger, fileCh, pkg, cleanedFileNames)
			if err != nil {
				_err = errorz.Wrapf(err, "parsing %v: %v", pkg.PkgPath, err)
			}
			p := filteringPackage{
				PkgPath:          pkg.PkgPath,
				Imports:          pkg.Imports,
				Directives:       directives,
				InlineDirectives: inlineDirectives,
			}
			collectedPackages[i] = p
		}(i, pkgs[i])
	}
	wg.Wait()
	close(fileCh)
	close(errCh)
	wg0.Wait()

	_err = errorz.Wrap(errorz.Combine(errs...), "can not parse packages")
	return
}

func cleanDir(cleanedFileNames map[string]bool, pkgDir string) error {
	dir, err := os.Open(pkgDir)
	if err != nil {
		return err
	}
	defer func() { must(dir.Close()) }()
	names, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, name := range names {
		if cleanedFileNames[name] {
			absFileName := filepath.Join(pkgDir, name)
			if err = os.Remove(absFileName); err != nil {
				return errorz.Wrapf(err, "can not remove file %v: %v", absFileName, err)
			}
		}
	}
	return nil
}

func parseDirectivesFromPackage(logger Logger, fileCh chan<- fileContent, pkg *packages.Package, cleanedFileNames map[string]bool) (directives, inlineDirectives []Directive, _err error) {
	for _, file := range pkg.CompiledGoFiles {
		if cleanedFileNames[filepath.Base(file)] {
			continue
		}
		body, err := os.ReadFile(file)
		if err != nil {
			return nil, nil, err
		}
		fileCh <- fileContent{Path: file, Body: body}

		errs := parseDirectivesFromBody(body, &directives, &inlineDirectives)
		if len(errs) != 0 {
			// ignore unknown directives
			for _, e := range errs {
				logger.Warn("ignored directive", "err", e)
			}
		}
	}
	return
}

var startDirective0 = []byte(startDirectiveStr0)
var startDirective1 = []byte(startDirectiveStr1)
var startDirective2 = []byte(startDirectiveStr2)

func maxStr(s string) string {
	idx := strings.IndexByte(s, '\n')
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func parseDirectivesFromBody(body []byte, directives, inlineDirectives *[]Directive) (errs []error) {

	// store processing directives
	var tmp []Directive
	lastIdx := -1
	if bytes.HasPrefix(body, startDirective0) || bytes.HasPrefix(body, startDirective1) || bytes.HasPrefix(body, startDirective2) {
		lastIdx = 0
	}
	for idx := 0; idx < len(body); idx++ {
		if body[idx] != '\n' {
			continue
		}

		// process the last found directive
		if lastIdx >= 0 {
			line := body[lastIdx:idx]
			lastIdx = -1

			directive, err := ParseDirective(string(line))
			if err != nil {
				errs = append(errs, err)
				continue
			}
			tmp = append(tmp, directive)
		}
		// directives are followed by a blank line, accept them
		if idx+1 < len(body) && body[idx+1] == '\n' {
			*directives = append(*directives, tmp...)
			tmp = tmp[:0]
		}
		// find the next directive
		if !bytes.HasPrefix(body[idx+1:], startDirective0) && !bytes.HasPrefix(body[idx+1:], startDirective1) && idx+1 != len(body) {
			// put directives not followed by a blank line into inline directives
			if inlineDirectives != nil {
				*inlineDirectives = append(*inlineDirectives, tmp...)
			}
			tmp = tmp[:0]
			continue
		}
		lastIdx = idx
	}
	// source file should end with a newline, so we don't process remaining lastIdx
	*directives = append(*directives, tmp...)
	return errs
}

func (ng *engine) validateConfig(cfg *Config) (_err error) {
	defer func() {
		if _err != nil {
			_err = errorz.Wrapf(_err, "config error: %v", _err)
		}
	}()

	// populate enabledPlugins
	if cfg.EnabledPlugins != nil {
		for name, enabled := range cfg.EnabledPlugins {
			if enabled {
				pl := ng.pluginsMap[name]
				if pl == nil {
					return errorz.Wrapf(nil, "plugin %v not found", name)
				}
				pl.enabled = true
				ng.enabledPlugins = append(ng.enabledPlugins, pl)
			}
		}
	} else {
		// enable all plugins
		for _, pl := range ng.plugins {
			pl.enabled = true
		}
		ng.enabledPlugins = ng.plugins
	}

	if cfg.GenerateFileName == nil {
		cfg.GenerateFileName = defaultFileNameGenerator(defaultGeneratedFileNameTpl)
	}

	if ng.bufPool.New == nil {
		ng.bufPool.New = func() any {
			return bytes.NewBuffer(make([]byte, 0, defaultBufSize))
		}
	}
	return nil
}

func (ng *engine) genFilename(input GenerateFileNameInput) string {
	return ng.xcfg.GenerateFileName(input)
}

func (ng *engine) writeFile(filePath string) (io.WriteCloser, error) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	ng.generatedFiles = append(ng.generatedFiles, filePath)
	return f, nil
}

func (ng *engine) execGoimport(files []string) error {
	var args []string
	args = append(args, ng.xcfg.GoimportsArgs...)
	args = append(args, "-w")
	args = append(args, files...)
	cmd := exec.Command("goimports", args...)
	ng.logger.Debug("goimports", "args", args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errorz.Wrapf(err, "goimports: %s\n\n%s\n", err, out)
	}
	return nil
}
