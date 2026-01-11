# Demo data for local development. Safe to re-run.

script = Script.find_or_create_by!(title: "Neon Nocturne") do |s|
  s.description = "A moody sci-fi noir where the city never sleeps."
end

script.ensure_vosd!

actors = ["Avery Quinn", "Jules Vega", "Riley Kade"].map do |name|
  Actor.find_or_create_by!(name: name)
end

characters = [
  { name: "Detective Marlowe", description: "Grizzled gumshoe with a synthetic heart.", actor: actors[0] },
  { name: "Nova", description: "Rebel courier with a hidden past.", actor: actors[1] },
  { name: "Archivist", description: "Keeper of the city memory vault.", actor: actors[2] }
].map do |attrs|
  Character.find_or_create_by!(script: script, name: attrs[:name]) do |c|
    c.description = attrs[:description]
    c.actor = attrs[:actor]
    c.role = "character"
  end
end

sections = [
  { name: "Opening Monologue", description: "Rain on neon glass.", ordering: 1 },
  { name: "The Alley Deal", description: "A trade goes sideways.", ordering: 2 },
  { name: "Closing Monologue", description: "A last cigarette.", ordering: 3 }
].map do |attrs|
  Section.find_or_create_by!(script: script, name: attrs[:name]) do |section|
    section.description = attrs[:description]
    section.ordering = attrs[:ordering]
  end
end

sections.each do |section|
  section.ensure_vosd!
end

props = [
  { name: "Trench Coat", description: "Weathered and soaked.", section: sections[0] },
  { name: "Encrypted Drive", description: "Glows faintly.", section: sections[1] },
  { name: "Old Lighter", description: "Refuses to die.", section: sections[2] }
]

props.each do |attrs|
  Prop.find_or_create_by!(name: attrs[:name], section: attrs[:section]) do |prop|
    prop.description = attrs[:description]
  end
end

sections.each_with_index do |section, index|
  characters.each_with_index do |character, char_index|
    next if index.zero? && char_index == 2

    CharacterSection.find_or_create_by!(section: section, character: character) do |cs|
      cs.on_stage = index != 1 || char_index != 0
    end
  end
end

writers = ["Alex", "Blair", "Casey"].map do |name|
  Writer.find_or_create_by!(name: name)
end

sections.each_with_index do |section, index|
  section.sub_sections.each do |sub_section|
    sub_section.update!(writer: writers[index % writers.length]) if sub_section.writer.nil?
  end
end

seed_line = sections.first.sub_sections.first.lines.first
if seed_line.nil?
  sections.first.sub_sections.first.lines.create!(
    character: script.VOSD,
    ordering: 0,
    text: "[The rain hammers the neon-lit street.]"
  )
end

demo_user = User.find_or_create_by!(email_address: "admin@example.com") do |user|
  user.password = "password"
  user.password_confirmation = "password"
  user.admin = true
end

UserToken.issue!(demo_user, name: "seed") if demo_user.user_tokens.empty?
