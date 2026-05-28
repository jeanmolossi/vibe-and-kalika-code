# Security Lens

Use this lens when the patch touches inputs, outputs, auth, storage, network, files, secrets, or permissions.

## What to check

- Inputs are validated before use.
- Sensitive values are not logged, echoed, cached, or serialized casually.
- Auth and authorization checks happen before privileged actions.
- File paths are constrained and cannot escape intended roots.
- SQL, shell, template, and expression execution paths are not injectable.
- External requests do not leak secrets or internal data.
- Error messages do not expose more than they should.

## High-risk categories

- Hardcoded secrets or tokens
- Missing authz checks
- Injection: SQL, shell, template, command, expression, prompt
- Path traversal and unsafe file reads/writes
- SSRF or unsafe outbound requests
- Insecure deserialization
- Weak randomness for security-sensitive operations
- Privacy leaks through logs or telemetry

## Questions to ask

- Can an attacker control this input?
- If yes, where is it validated or escaped?
- What happens on failure?
- Does the error path leak secrets or bypass checks?
- Does this change widen the blast radius of a compromised input?

## Red flags

- `eval`, `exec`, shell interpolation, or raw query concatenation
- Logging tokens, passwords, session IDs, personal data, or internal URLs
- Trusting client-side checks for server-side authorization
- File operations built from unchecked user paths
- Security logic hidden inside helper code that reviewers may miss

## Review rule

Security issues are usually blockers unless there is a very strong mitigation already in place and clearly visible in the patch.
