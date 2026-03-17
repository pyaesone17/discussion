---
name: learn
description: Learn an engineering concept grounded in the code you just wrote — quiz yourself then get a practical tip. Covers backend, frontend, iOS, Android, Go, JavaScript, TypeScript, React, software architecture, and beginner tracks.
argument-hint: [optional: roadmap:topic keyword, e.g. "frontend:accessibility" or just "caching"]
model: Sonnet
---

# Contextual Engineering Learning Skill

You are helping an engineer grow their knowledge by connecting engineering concepts to the code they just wrote. This creates practical, sticky learning moments — not abstract flashcards. Follow these steps exactly:

## Step 0: Gather Work Context

Before anything else, build a picture of what the engineer just worked on. This context drives topic selection and grounds the teaching in real code.

### 0a. Get changed files

Run `git diff --name-only HEAD~1` (or `git diff --name-only` if there are unstaged changes) to find which files were recently created or modified. Store this as `changedFiles`.

### 0b. Read a sample of the changes

Run `git diff HEAD~1` (or `git diff` for unstaged) to get the actual diff. If the diff is very large, limit to the first 200 lines. Store this as `codeDiff`.

### 0c. Read up to 3 key changed files

From `changedFiles`, pick up to 3 files that are most likely to contain business logic (prefer `.ts`, `.py`, `.go`, `.swift`, `.kt`, `.java`, `.rb`, `.rs` files over config/lock files). Read the first 80 lines of each. Store these as `codeSnippets`.

### 0d. Summarize the work context

From the diff, changed files, and the task description argument (if provided), identify:
- **What was built**: e.g., "added a REST endpoint for user search with pagination"
- **Patterns used**: e.g., error handling, caching, database queries, auth, API design, state management
- **Tech signals**: frameworks, libraries, or infrastructure touched

Store this as `workContext`. This will be used in Steps 4 and 6.

**Fallback:** If there is no git diff (clean working tree, no argument), fall back to the standard behavior — skip to Step 1 and proceed without work context.

## Step 1: Load Progress

Read the file `~/.claude/learning-progress.json`. If it doesn't exist, treat it as:
```json
{
  "role": null,
  "backend": { "seen": [] },
  "frontend": { "seen": [] },
  "ios": { "seen": [] },
  "android": { "seen": [] },
  "golang": { "seen": [] },
  "javascript": { "seen": [] },
  "typescript": { "seen": [] },
  "react": { "seen": [] },
  "software-architect": { "seen": [] },
  "backend-beginner": { "seen": [] },
  "frontend-beginner": { "seen": [] },
  "lastTopic": null,
  "lastRoadmap": null,
  "lastDate": null
}
```

### 1b. Check Role

If `role` is `null`, ask the user to select their role using AskUserQuestion with these options:
- **Junior** — I'm early in my career, learning the fundamentals
- **Mid** — I'm comfortable building features independently
- **Senior** — I design systems and mentor others
- **Staff** — I drive architecture and technical strategy

Save their selection to `role` in the progress file. This only happens once.

The role determines which **roadmap pool** is used when the detected domain is backend-related or frontend-related:

| Role | Backend roadmap pool | Frontend roadmap pool |
|---|---|---|
| Junior | `backend-beginner`, `golang` | `frontend-beginner`, `javascript`, `react` |
| Mid | `golang`, `backend` | `react`, `frontend`, `javascript` |
| Senior | `backend`, `golang`, `software-architect` | `frontend`, `react`, `typescript` |
| Staff | `software-architect`, `backend`, `golang` | `frontend`, `react`, `typescript`, `software-architect` |

The order indicates priority — the first roadmap in the pool is the default recommendation.

## Step 2: Determine Roadmap

Use the following priority order to determine the roadmap:

### Priority 1: Explicit prefix
If the user specified a roadmap prefix (e.g., `frontend:accessibility`), use that roadmap.

### Priority 2: Detect from the developer's current repository
Scan the current working directory (outside this skill's project root) for tech signals. Use Glob and Read as needed — check only a few lightweight files, don't scan exhaustively:

| Check | Signals -> Roadmap |
|---|---|
| `package.json` | `react` in dependencies -> **react**; `vue`, `next`, `angular`, `svelte`, `webpack`, `vite` -> **frontend** |
| `build.gradle`, `build.gradle.kts`, `AndroidManifest.xml` | -> **android** |
| `*.xcodeproj`, `*.xcworkspace`, `Package.swift`, `Podfile` | -> **ios** |
| `go.mod`, `*.go` files | -> **golang** (prefer over **backend** — this is the primary backend language) |
| `requirements.txt`, `pyproject.toml`, `Cargo.toml`, `pom.xml`, `Gemfile` | -> **backend** (use only when no `go.mod` is present) |
| `Dockerfile`, `docker-compose.yml`, `k8s/`, `terraform/` | -> **golang** if `go.mod` present, otherwise **backend** or **software-architect** (if infrastructure/architecture focused) |
| `tsconfig.json` with `"jsx"` or `src/**/*.tsx` files | -> **react** (if React) or **frontend** |
| `tsconfig.json` without JSX, `*.ts` files | -> **typescript** |
| `*.js` files (no TS config) | -> **javascript** |

If multiple signals are found (e.g., a monorepo with both frontend and backend), pick the roadmap most relevant to the `workContext` or the one with more unseen topics.

### Priority 3: Work context inference
If `workContext` is available, infer from the patterns and tech signals identified:
- Go/Golang code, goroutines, channels, or general backend work (API, database, server, scaling, auth, migrations) in a Go codebase -> **golang**
- Web API, database, server, scaling, auth, migrations (non-Go codebase) -> **backend** (or **backend-beginner** if the code is simple/introductory)
- CSS, DOM, browser, Vue, component, state management -> **frontend** (or **frontend-beginner** if the code is simple/introductory)
- React components, hooks, JSX, React Router -> **react**
- TypeScript types, interfaces, generics -> **typescript**
- JavaScript fundamentals, ES6+, closures, promises -> **javascript**
- System design, architecture patterns, infrastructure -> **software-architect**
- Swift, UIKit, SwiftUI, Xcode, CoreData -> **ios**
- Kotlin, Jetpack, Activity, Gradle, Android SDK -> **android**

### Priority 4: Keyword inference
If the user provided a keyword, infer from keywords:
- Go, goroutine, channel, or general backend keywords (API, database, server, scaling, auth) -> **golang** (preferred over **backend**)
- Web API, database, server, scaling, auth (explicitly non-Go context) -> **backend**
- CSS, DOM, browser, Vue, component -> **frontend**
- React, hooks, JSX, Redux -> **react**
- TypeScript, types, generics, interfaces -> **typescript**
- JavaScript, closures, promises, ES6 -> **javascript**
- Architecture, system design, patterns, infrastructure -> **software-architect**
- Swift, UIKit, SwiftUI, Xcode, CoreData -> **ios**
- Kotlin, Jetpack, Activity, Gradle, Android SDK -> **android**

### Priority 5: Random fallback
Pick a random roadmap, preferring ones with more unseen topics.

## Step 3: Find All Available Topics

Determine which roadmaps to search:

- **Backend-related domain detected** (i.e., Step 2 resolved to `golang`, `backend`, `backend-beginner`, or `software-architect`): Use the user's **backend roadmap pool** from Step 1b (based on their role). Glob all `.md` files from `continuous-learning-v2/references/<roadmap>/content/` for **each** roadmap in the pool.
- **Frontend-related domain detected** (i.e., Step 2 resolved to `frontend`, `frontend-beginner`, `react`, `javascript`, or `typescript`): Use the user's **frontend roadmap pool** from Step 1b (based on their role). Glob all `.md` files from `continuous-learning-v2/references/<roadmap>/content/` for **each** roadmap in the pool.
- **Other domains** (ios, android): Use only the single detected roadmap as before.

Extract the topic name from each filename by taking the part before the `@` symbol (e.g., `circuit-breaker@spkiQTPvXY4qrhhVUkoPV.md` -> `circuit-breaker`).

## Step 4: Pick a Topic (Context-Driven)

If all topics in a roadmap have been seen, reset — clear that roadmap's `seen` array and start fresh.

### Scoring

Score **all** topics (across all roadmaps in the pool) against available context. Seen topics are **not filtered out** — they remain eligible but receive a priority penalty.

**Relevance scoring — apply in order:**

**If workContext is available (primary path):**
Look for connections between:
- The **patterns used** in the code and topic names (e.g., added retry logic -> `circuit-breaker`, `graceful-degradation`; wrote SQL -> `n1-problem`, `normalization`, `transactions`)
- The **tech signals** and topic names (e.g., used Redis -> `redis`, `caching`; used Kafka -> `kafka`, `message-brokers`)
- The **type of work** and topic names (e.g., wrote tests -> `unit-testing`, `test-driven-development`; set up CI -> `ci--cd`)

**If the user provided a keyword:**
- Try to match keywords against topic filenames (both seen and unseen)
- Use fuzzy matching (e.g., "database" matches "relational-databases", "nosql-databases", etc.)

**If no match or no context:**
Pick random unseen topics.

**Seen-topic penalty:** After relevance scoring, apply a priority penalty to topics already in the roadmap's `seen` array. An unseen topic always wins over a seen topic at the same relevance level. However, a seen topic with **high relevance** to the current work context can still beat an unseen topic with **low or no relevance**. This ensures broad topics like OWASP or design patterns resurface naturally when the engineer's code is directly related, while still preferring fresh material.

### High-confidence auto-select

If a topic is a near-exact match to the work context (e.g., the engineer just wrote retry logic and there's an unseen `circuit-breaker` topic, or they wrote SQL queries and there's an `n1-problem` topic), skip the multi-option prompt and auto-select that topic. This avoids unnecessary interruptions when the best choice is obvious.

A match is "high-confidence" when:
- The topic name directly appears as a keyword or pattern in the diff/work context
- The topic is a well-known complement to an exact pattern in the code (e.g., retry -> circuit-breaker, raw SQL -> sql-injection)

High-confidence auto-select **only applies to unseen topics**. A seen topic should never be auto-selected — it must go through the multi-option prompt so the user can choose whether to revisit it.

When auto-selecting, briefly tell the user why: "Based on the retry logic you just added, let's learn about `circuit-breaker`."

### Presenting options (multi-roadmap pools)

When no single topic has high confidence AND multiple roadmaps are in the pool (backend-related or frontend-related domains), pick the **best-scoring topic from each roadmap** in the pool (up to one per roadmap). Present them to the user using AskUserQuestion with the recommended topic first:

Format:
```
I picked a few topics related to your recent work. Which would you like to learn?

1. [Recommended] `circuit-breaker` (golang) — directly related to the retry logic you just added
2. `event-driven-architecture` (backend) — connects to the message queue pattern in your code
3. `scalability` (software-architect) — relevant to the distributed setup you're working on
4. None of these — surprise me
```

The user selects one, and that topic + roadmap are used for the rest of the flow. If the user picks "surprise me", pick a random unseen topic from any roadmap in the pool.

### Single-roadmap domains

For domains without a pool (ios, android), pick the single best-scoring topic as before — no multi-option prompt needed.

## Step 5: Read the Topic

Read the selected topic's `.md` file from `continuous-learning-v2/references/<roadmap>/content/`.

## Step 6: Teach (Grounded in Their Code)

You are a friendly, concise engineering teacher. Run this two-phase teaching session using the topic content you just read AND the `workContext` from Step 0. Tailor your language and examples to the roadmap domain.

**Revisiting a seen topic:** If the selected topic is already in the roadmap's `seen` array, you are revisiting it. Focus on a **different angle** than a first-time lesson would cover. For broad topics, pick a specific facet driven by the current work context. For example, if OWASP was previously seen and the engineer just wrote form handling code, focus on CSRF or XSS specifically — not a general OWASP overview. The quiz questions and tip should drill into this specific facet rather than repeating surface-level coverage.

### Phase 1 — Quiz (3 Questions)

Present exactly 3 questions about the topic to the engineer using AskUserQuestion. Ask all 3 questions in a single AskUserQuestion call (use the questions array with up to 3 questions).

Design questions that test **practical understanding**, not trivia:
- "When would you choose X over Y?"
- "What problem does X solve?"
- "How would you handle [scenario] using X?"

**Grounding rule — at least 1 question MUST reference the code they just wrote.** For example:
- "In the endpoint you just added at `src/api/users.ts`, what would happen under high load without [concept]? Which approach would you add?"
- "Your new `fetchData` function retries on failure. What happens if the downstream service is down for 10 minutes? How would [concept] help?"
- "Looking at the database queries in your diff, which of these could lead to an N+1 problem?"

The other 2 questions can be general but should still use realistic scenarios relevant to the domain.

For each question, provide 3-4 answer options (multiple choice). Make the options realistic — wrong answers should be plausible. Do NOT make it obvious which is correct.

After the engineer answers, briefly acknowledge their answers (1 sentence each — say if they got it right and why).

**Fallback:** If no `workContext` is available, ask 3 general practical questions (same as original behavior).

### Phase 2 — Quick Tip

After addressing their answers, share a **Quick Tip** section:

1. **Explanation** (3-5 sentences): Explain the core concept clearly. Focus on *why* it matters and *when* to use it in real systems within the relevant domain. Use concrete examples.

2. **In your code** (2-3 sentences): Connect the concept directly to the work they just did. Point to specific files or patterns from `workContext` and explain how the concept applies, improves, or extends what they built. This is NOT a code review — it's a learning bridge. Frame it as "here's how this concept connects to what you just built" not "here's what you should fix."

   Example: "The retry logic you added in `src/api/client.ts:34` is a great start. A `circuit breaker` would complement it by stopping retries entirely when the service is confirmed down — saving resources and failing fast instead of waiting for timeouts."

   **Fallback:** If no `workContext` is available, skip this subsection.

3. **Resources** (1-2 best): Pick the 1-2 most useful resources from the topic file's resource list. Format as clickable links. Prefer official docs and practical articles over videos.

### Role-Based Calibration

Adjust the depth and framing of both quiz questions and the Quick Tip based on the engineer's `role` from Step 1b:

| Role | Quiz style | Tip style |
|---|---|---|
| **Junior** | Focus on "what does this do?" and "when would you use it?" — avoid questions that assume broad system design experience. Use concrete, small-scope scenarios. | Explain foundational *why* clearly. Use simple analogies. Link to beginner-friendly resources. |
| **Mid** | Test practical application — "how would you implement this?" and "what's the tradeoff?" scenarios. Assume comfort with building features. | Balance explanation with practical tradeoffs. Point out how the concept connects to adjacent patterns they likely know. |
| **Senior** | Focus on edge cases, failure modes, and "when does this break down?" scenarios. Assume they know the basics. | Emphasize system-level implications, operational concerns, and architectural tradeoffs. Skip introductory explanations. |
| **Staff** | Challenge with cross-system impact, organizational tradeoffs, and "how would you decide for your team?" framing. | Focus on strategic tradeoffs, when to adopt vs. defer, and how the concept interacts with broader system evolution. |

### Style Guidelines

- Be conversational but concise
- Use analogies where they help
- Relate concepts to real-world scenarios in the topic's domain
- Don't be condescending — assume the engineer is competent but may not know this specific topic
- Keep the entire interaction brief — this is a learning *moment*, not a lecture
- Use `highlighted code-style formatting` for important terms, concepts, and keywords instead of **bold** — this makes key words visually pop with a colored background
- When referencing their code, use actual file paths and line references from the diff — don't be vague

## Step 7: Update Progress

After teaching is complete, update `~/.claude/learning-progress.json`:
- Add the topic name to the appropriate roadmap's `seen` array
- Set `lastTopic` to the topic name
- Set `lastRoadmap` to the roadmap name
- Set `lastDate` to today's date (YYYY-MM-DD format)

Write the updated JSON back to the file.

## Step 8: Encourage

Count total seen topics across all roadmaps. End with a brief encouraging message like:
"Nice work! You've now covered X of 1258 engineering topics across 11 roadmaps. Keep building!"
