## 📖 Description
<!-- Clearly describe what this PR does -->
- Update contract service to query matchbox for selected policy information  
- Jira: https://jira.kdc.capitalone.com/browse/CT2168T-41

---

## 🎯 Type of Change
<!-- Mark relevant option with x -->
- [ ] 🐞 Bug fix
- [ ] ✨ New feature
- [ ] ⚡ Performance improvement
- [ ] ♻️ Refactoring
- [ ] 🧪 Test update
- [ ] 📝 Documentation update

---

## 🔗 Related Issue / Jira
- Jira ID: **CT2168T-41**
- Github Issue (if any): #

---

## 🧪 How Has This Been Tested?
<!-- Explain how you tested your changes -->
- [ ] Unit Tests
- [ ] Integration Tests
- [ ] Manual Testing
- [ ] Tested in Dev Environment

**Test Details:**
<!-- Example -->
- Verified API response from matchbox service
- Checked fallback scenarios

---

## ✅ Acceptance Criteria
<!-- Ensure all are checked before raising PR -->
- [ ] Code compiles successfully
- [ ] Unit tests added/updated (Coverage ≥ 90%)
- [ ] Functional tests updated (if applicable)
- [ ] Documentation updated (if applicable)
- [ ] No breaking changes OR properly documented
- [ ] Pipeline build is passing

---

## 📸 Screenshots / Logs (if applicable)
<!-- Add screenshots, logs, or API responses -->

---

## ⚠️ Risks & Impact
<!-- Mention any risks -->
- Impacted Services:
- Backward Compatibility:
- Migration Required:

---

## 🚀 Deployment Notes
<!-- Any special deployment steps -->
- No special steps required / OR mention steps

---

## 👀 Reviewer Notes
<!-- Anything specific you want reviewer to focus on -->
- Please review matchbox integration logic
- Focus on retry handling

---

## 📋 Checklist
- [ ] PR title follows standard: `<TeamName> | <JiraID> | <Title>`
- [ ] Code follows project standards
- [ ] Self-review completed
- [ ] No unnecessary logs/debug code
