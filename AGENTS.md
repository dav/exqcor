# AGENTS.md — Rails 2 -> Rails 8 Upgrade

## Goals
- Analyze the current Rails 3.x era project to understand how it was implemented. See below for a description of the project and how it was used.
- Create a Rails 8.x version of this application

## Understanding the exqcor project, basic idea

This system was used in stage productions that followed a particular format. The script for the stage production was written by the audience the night of the show, and cold-read performed by improv actors.

Each production had a theme (Soap Opera, Film Noir, SciFi. Upon arrival at the theater, the audience would join a gathering on the stage where food and drinks where available. They would mingle with the actors for nights performance. The actors would be in costume and in character (a character compatible with the theme, e.g. Startship Captain, Rebel Commander, Space Smuggler, etc). Each audience member would get a number. After some time getting to observe the actors in character the audience members would start being called by number into the back stage area, where three computers were set up, one for Act 1, Act2 and Act 3.

When sat down at the computer, the audience member (now "writer") would see a UI that allowed them to select a character or VOSD (stage direction) and enter dialogue/direction. They would be given 5 minutes to write, and when the timer went off, they got up and the next write would sit down to add to the script for that Act. Crucially, the next writer would only see the last line written by the previous writer (hence "exquisite corpse"). After enough lines were written the script for that act could be printed and performed on stage minutes later.

## Previous codex analysis of project

Current data structures

Play: title/description; owns characters, sections, and provides helper methods like VOSD autogeneration (Voice of Stage Directions) plus duplication logic that copies sections and characters for templates (app/models/play.rb; sed -n '1,200p' app/models/play.rb).

Section (scene): belongs to a play, has props, sub_sections, and a character_sections join that tracks which characters are available/on stage; auto-creates a first sub_section and ensures VOSD is included (app/models/section.rb; sed -n '1,200p' app/models/section.rb).

SubSection (writer turn): ordered slices inside a section; contains ordered lines and can calculate writing duration between first/last non-priming lines (app/models/sub_section.rb; sed -n '1,200p' app/models/sub_section.rb).

Line: text plus ordering, ties to a character and a sub_section (app/models/line.rb; sed -n '1,200p' app/models/line.rb).

Character: name/description, belongs to a play and optionally an actor; uniqueness scoped per play; protects the VOSD record from deletion (app/models/character.rb; sed -n '1,200p' app/models/character.rb).

Actor: simple name model linked to characters (app/models/actor.rb; sed -n '1,200p' app/models/actor.rb).

Prop: name/description tied to a section (app/models/prop.rb; sed -n '1,200p' app/models/prop.rb).
Join table character_sections adds on_stage boolean (db/schema.rb; sed -n '1,200p' db/schema.rb).
Key workflows and functionality

Admin play setup: Create play, add characters (including automatic VOSD), assign actors, and build scenes/sections. UI uses checkboxes to mark which characters are available/on-stage per section and to manage props (app/views/plays/_form.html.erb, app/views/plays/_sections.html.erb; commands: sed -n '1,240p' app/views/plays/_form.html.erb, sed -n '1,240p' app/views/plays/_sections.html.erb).

Timed writing sessions: Each section contains multiple writer “turns” (sub_sections). Writers see only the previous turn’s last line, select a character, and add lines through a Backbone-powered interface with a 5-minute countdown. After time expires, they’re prompted to hand off to the next writer (app/views/sub_sections/edit.html.erb, app/assets/javascripts/routers/sub_sections_router.js.coffee, app/assets/javascripts/views/sub_sections/sub_sections_writer.js.coffee, app/assets/templates/sub_sections/writer.jst.eco; commands: sed -n '1,200p' app/views/sub_sections/edit.html.erb, sed -n '1,200p' app/assets/javascripts/routers/sub_sections_router.js.coffee, sed -n '1,240p' app/assets/javascripts/views/sub_sections/sub_sections_writer.js.coffee, sed -n '1,200p' app/assets/templates/sub_sections/writer.jst.eco).

Script assembly & printing: Full scripts can be generated per play, section, actor, or character, highlighting relevant lines and marking stage directions (VOSD). There’s also a “script package” that bundles actor-specific scenes (app/views/plays/full_script.html.erb, app/views/sections/_section_actor_script.html.erb, app/views/sections/_section_character_script.html.erb, app/views/sections/script_package.html.erb, app/views/sections/_section_script_header.html.erb; commands: sed -n '1,200p' app/views/plays/full_script.html.erb, sed -n '1,200p' app/views/sections/_section_actor_script.html.erb, sed -n '1,200p' app/views/sections/_section_character_script.html.erb, sed -n '1,200p' app/views/sections/script_package.html.erb, sed -n '1,200p' app/views/sections/_section_script_header.html.erb).
Statistics/reporting: Simple averages on lines per writer/scene and writing duration per scene/turn (app/views/plays/_report.html.erb; sed -n '1,200p' app/views/plays/_report.html.erb).

Routing/API shape: Nested REST routes (config/routes.rb; sed -n '1,200p' config/routes.rb) with Backbone models/collections overriding URLs for create/delete, meaning the mobile/API layer would need equivalents.

## Constraints
- Prefer the smallest viable change.
- Don’t introduce new dependencies unless required.
- Document any removed/renamed features.

## Workflow
1. Create for ruby 3.4 and Rails 8
2. Create the new rails app in a `modern/` sub directory with its own Gemfile there
3. After each step: bundle install, run tests, boot app, fix deprecations.
4. Keep a running UPGRADE_NOTES.md with decisions and breakages.

## Commands
- Use `bin/rails` when present.
- Never run destructive commands without asking (db drops, etc.).
