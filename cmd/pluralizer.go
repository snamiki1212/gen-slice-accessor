package cmd

import "regexp"

type pluralizer struct {
	rules pluralizeRules
}

// REF: https://github.com/gertd/go-pluralize/blob/master/pluralize.go#L317
type (
	pluralizeRule struct {
		expression  *regexp.Regexp
		replacement string
	}
	pluralizeRules []pluralizeRule
	strRule        struct {
		expr string
		rep  string
	}
	strRules []strRule
)

func (sc strRules) toPluralizeRules() []pluralizeRule {
	rules := make([]pluralizeRule, 0, len(sc))
	for _, v := range sc {
		rules = append(rules, v.toPluralizeRule())
	}
	return rules
}

func (s strRule) toPluralizeRule() pluralizeRule {
	return pluralizeRule{newRegexpRule(s.expr), s.rep}
}

func (p pluralizer) replace(word string, rules []pluralizeRule) string {
	for _, r := range rules {
		if r.expression.MatchString(word) {
			return p.doReplace(word, r)
		}
	}
	return word
}

func (p pluralizer) doReplace(word string, rule pluralizeRule) string {
	return rule.expression.ReplaceAllString(word, rule.replacement)
}

// NOTE: (?i) is case insentive flag(https://stackoverflow.com/questions/15326421/how-do-i-do-a-case-insensitive-regular-expression-in-go)
func newDefaultStrRules() strRules {
	return strRules{
		{`(?i).*`, `${0}s`},
	}
}

func isExpr(s string) bool {
	return s[:1] == `(`
}

func newRegexpRule(rule string) *regexp.Regexp {
	rl := func() string {
		if isExpr(rule) {
			return rule
		}
		return `(?i)^` + rule + `$`
	}()
	return regexp.MustCompile(rl)
}

func newPluralizer(rulesList ...strRules) pluralizer {
	rules := make(strRules, 0)
	if len(rulesList) > 0 {
		rs := rulesList[0] // only get the first elements
		rules = append(rules, rs...)
	}
	rules = append(rules, newDefaultStrRules()...)
	return pluralizer{rules.toPluralizeRules()}
}

func (p pluralizer) pluralize(str string) string {
	return p.replace(str, p.rules)
}
