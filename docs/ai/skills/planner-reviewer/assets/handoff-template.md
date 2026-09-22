# Implementation Handoff — {{task title}}

- Task ID / instruction revision: {{ID / revision}}
- Date and environment: {{date, timezone, relevant runtime}}
- Code revision: {{commit hash; if uncommitted, base hash and dirty state}}
- Implementation status: {{IMPLEMENTED / PARTIAL / BLOCKED}}
- Review status: PENDING REVIEW

## Result

{{What now works, what remains incomplete and the impact on the user.
Do not declare the task accepted on behalf of the reviewer.}}

## Changes

{{Group meaningful changes by behavior; include relevant file paths.
Mention schema/API/config changes and setup steps only when applicable.}}

## Acceptance evidence

| AC ID | Status: PASS / FAIL / NOT VERIFIED | Evidence and limitation |
| --- | --- | --- |
| AC-01 | {{status}} | {{Check/result, date, code revision, artifact or source location}} |

PASS requires evidence appropriate to the criterion. Code inspection alone cannot
establish successful live behavior. Keep failures and unexecuted checks visible;
identify any exceptions the user explicitly accepted without changing their status.

## Checks actually performed

| Check or command | Result | Evidence / relevant environment |
| --- | --- | --- |
| {{Check}} | {{PASS / FAIL / NOT VERIFIED}} | {{Actual result, counts or failure details}} |

{{Clearly label historical results, current results, manual user reports and checks
not run. Never infer success from silence or claim tests on a different revision
validate the current one without explaining their applicability.}}

## Decisions and deviations

{{Implementation choices that matter to review; differences from the brief,
reason, impact and whether the user approved any scope change. State none if none.}}

## Remaining issues and next action

{{Reproduction steps, impact, blockers, missing evidence and specific next check.
Separate required fixes from optional improvements. Do not propose the next feature
as a replacement for finishing this task.}}

## Working tree and runtime

{{Branch, commits created, push status, uncommitted changes; services left running,
test data created/removed and other meaningful side effects. No secrets.}}

## Review pointers

{{Smallest set of files, diffs and evidence the planner needs to review the result.}}
