# Coding Instructions — {{task title}}

- Task ID: {{stable ID}}
- Revision: {{r1 / r2 / ...}}
- Date: {{date and timezone}}
- Kind: {{implementation / correction / evidence-only}}
- Readiness: {{DRAFT / READY}}
- Supersedes: {{previous brief or none}}

## Outcome and context

{{User-visible outcome; current implementation facts with sources/dates; dependencies.
Include enough context to execute without access to the planning conversation.}}

## Agreed decisions

{{Final product rules, units/defaults/bounds where agreed, and constraints.
Separate assumptions from confirmed decisions.}}

## Scope

{{Required work and explicit exclusions. Preserve unrelated user changes.}}

## Repository guidance

{{Relevant repository instructions and confirmed source locations.
If unavailable, direct the coding agent to discover them; do not invent file paths.
Follow established conventions. Flag material conflicts before dependent changes.}}

## Acceptance criteria

| ID | Required observable behavior | Expected evidence |
| --- | --- | --- |
| AC-01 | {{Concrete input/action and expected result}} | {{Appropriate check}} |

## Verification

{{Relevant commands if known, manual scenarios and expected outcomes.
Specify isolation/cleanup when test data is necessary. Do not expose secrets.
Record unrun checks as NOT VERIFIED and explain why.}}

## Correction context

{{For a correction or evidence request: previous handoff/revision, observed issue,
affected AC IDs, evidence, permitted fix scope and checks to rerun.
Omit this section for a first implementation.}}

## Open decisions and blockers

{{State none or list genuinely unresolved items. Required product decisions keep
this brief DRAFT. An environmental blocker is not evidence of a code defect.}}

## Delivery

Produce `{{task-id}}-handoff-{{revision}}.md` using the handoff format supplied below
or alongside this brief. Report actual changes and evidence, not just intentions.
State commit/push status; perform those actions only when authorized by the user.
Follow the repository's context-update rules. Do not start the next feature.

{{Insert the complete handoff template here, or link the accompanying accessible file.}}
