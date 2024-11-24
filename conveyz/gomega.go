package conveyz

import (
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func Equal(expected any) types.GomegaMatcher {
	return gomega.Equal(expected)
}
func BeEquivalentTo(expected any) types.GomegaMatcher {
	return gomega.BeEquivalentTo(expected)
}
func BeComparableTo(expected any, opts ...cmp.Option) types.GomegaMatcher {
	return gomega.BeComparableTo(expected, opts...)
}
func BeIdenticalTo(expected any) types.GomegaMatcher {
	return gomega.BeIdenticalTo(expected)
}
func BeNil() types.GomegaMatcher {
	return gomega.BeNil()
}
func BeTrue() types.GomegaMatcher {
	return gomega.BeTrue()
}
func BeFalse() types.GomegaMatcher {
	return gomega.BeFalse()
}
func BeTrueBecause(format string, args ...any) types.GomegaMatcher {
	return gomega.BeTrueBecause(format, args...)
}
func BeFalseBecause(format string, args ...any) types.GomegaMatcher {
	return gomega.BeFalseBecause(format, args...)
}
func HaveOccurred() types.GomegaMatcher {
	return gomega.HaveOccurred()
}
func Succeed() types.GomegaMatcher {
	return gomega.Succeed()
}
func MatchError(expected any, functionErrorDescription ...any) types.GomegaMatcher {
	return gomega.MatchError(expected, functionErrorDescription...)
}
func BeClosed() types.GomegaMatcher {
	return gomega.BeClosed()
}
func Receive(args ...any) types.GomegaMatcher {
	return gomega.Receive(args...)
}
func BeSent(arg any) types.GomegaMatcher {
	return gomega.BeSent(arg)
}
func MatchRegexp(regexp string, args ...any) types.GomegaMatcher {
	return gomega.MatchRegexp(regexp, args...)
}
func ContainSubstring(substr string, args ...any) types.GomegaMatcher {
	return gomega.ContainSubstring(substr, args...)
}
func HavePrefix(prefix string, args ...any) types.GomegaMatcher {
	return gomega.HavePrefix(prefix, args...)
}
func HaveSuffix(suffix string, args ...any) types.GomegaMatcher {
	return gomega.HaveSuffix(suffix, args...)
}
func MatchJSON(json any) types.GomegaMatcher {
	return gomega.MatchJSON(json)
}
func MatchXML(xml any) types.GomegaMatcher {
	return gomega.MatchXML(xml)
}
func MatchYAML(yaml any) types.GomegaMatcher {
	return gomega.MatchYAML(yaml)
}
func BeEmpty() types.GomegaMatcher {
	return gomega.BeEmpty()
}
func HaveLen(count int) types.GomegaMatcher {
	return gomega.HaveLen(count)
}
func HaveCap(count int) types.GomegaMatcher {
	return gomega.HaveCap(count)
}
func BeZero() types.GomegaMatcher {
	return gomega.BeZero()
}
func ContainElement(element any, result ...any) types.GomegaMatcher {
	return gomega.ContainElement(element, result...)
}
func BeElementOf(elements ...any) types.GomegaMatcher {
	return gomega.BeElementOf(elements...)
}
func BeKeyOf(element any) types.GomegaMatcher {
	return gomega.BeKeyOf(element)
}
func ConsistOf(elements ...any) types.GomegaMatcher {
	return gomega.ConsistOf(elements...)
}
func HaveExactElements(elements ...any) types.GomegaMatcher {
	return gomega.HaveExactElements(elements...)
}
func ContainElements(elements ...any) types.GomegaMatcher {
	return gomega.ContainElements(elements...)
}
func HaveEach(element any) types.GomegaMatcher {
	return gomega.HaveEach(element)
}
func HaveKey(key any) types.GomegaMatcher {
	return gomega.HaveKey(key)
}
func HaveKeyWithValue(key any, value any) types.GomegaMatcher {
	return gomega.HaveKeyWithValue(key, value)
}
func HaveField(field string, expected any) types.GomegaMatcher {
	return gomega.HaveField(field, expected)
}
func HaveExistingField(field string) types.GomegaMatcher {
	return gomega.HaveExistingField(field)
}
func HaveValue(matcher types.GomegaMatcher) types.GomegaMatcher {
	return gomega.HaveValue(matcher)
}
func BeNumerically(comparator string, compareTo ...any) types.GomegaMatcher {
	return gomega.BeNumerically(comparator, compareTo...)
}
func BeTemporally(comparator string, compareTo time.Time, threshold ...time.Duration) types.GomegaMatcher {
	return gomega.BeTemporally(comparator, compareTo, threshold...)
}
func BeAssignableToTypeOf(expected any) types.GomegaMatcher {
	return gomega.BeAssignableToTypeOf(expected)
}
func Panic() types.GomegaMatcher {
	return gomega.Panic()
}
func PanicWith(expected any) types.GomegaMatcher {
	return gomega.PanicWith(expected)
}
func BeAnExistingFile() types.GomegaMatcher {
	return gomega.BeAnExistingFile()
}
func BeARegularFile() types.GomegaMatcher {
	return gomega.BeARegularFile()
}
func BeADirectory() types.GomegaMatcher {
	return gomega.BeADirectory()
}
func HaveHTTPStatus(expected ...any) types.GomegaMatcher {
	return gomega.HaveHTTPStatus(expected...)
}
func HaveHTTPHeaderWithValue(header string, value any) types.GomegaMatcher {
	return gomega.HaveHTTPHeaderWithValue(header, value)
}
func HaveHTTPBody(expected any) types.GomegaMatcher {
	return gomega.HaveHTTPBody(expected)
}
func And(ms ...types.GomegaMatcher) types.GomegaMatcher {
	return gomega.And(ms...)
}
func SatisfyAll(matchers ...types.GomegaMatcher) types.GomegaMatcher {
	return gomega.SatisfyAll(matchers...)
}
func Or(ms ...types.GomegaMatcher) types.GomegaMatcher {
	return gomega.Or(ms...)
}
func SatisfyAny(matchers ...types.GomegaMatcher) types.GomegaMatcher {
	return gomega.SatisfyAny(matchers...)
}
func Not(matcher types.GomegaMatcher) types.GomegaMatcher {
	return gomega.Not(matcher)
}
func WithTransform(transform any, matcher types.GomegaMatcher) types.GomegaMatcher {
	return gomega.WithTransform(transform, matcher)
}
func Satisfy(predicate any) types.GomegaMatcher {
	return gomega.Satisfy(predicate)
}
