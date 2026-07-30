CREATE TABLE scripts (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    theme TEXT NOT NULL DEFAULT '',
    writing_seconds INTEGER NOT NULL DEFAULT 300,
    station_mode TEXT NOT NULL DEFAULT 'station' CHECK (station_mode IN ('station', 'byod')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE actors (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    bio TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE characters (
    id INTEGER PRIMARY KEY,
    script_id INTEGER NOT NULL REFERENCES scripts(id) ON DELETE CASCADE,
    actor_id INTEGER REFERENCES actors(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT 'character' CHECK (role IN ('character', 'vosd')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE (script_id, name)
);

CREATE UNIQUE INDEX idx_characters_one_vosd ON characters(script_id) WHERE role = 'vosd';

CREATE TABLE sections (
    id INTEGER PRIMARY KEY,
    script_id INTEGER NOT NULL REFERENCES scripts(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    ordering INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'writing', 'complete')),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE character_sections (
    id INTEGER PRIMARY KEY,
    character_id INTEGER NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    section_id INTEGER NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    on_stage INTEGER NOT NULL DEFAULT 0,
    UNIQUE (character_id, section_id)
);

CREATE TABLE props (
    id INTEGER PRIMARY KEY,
    section_id INTEGER NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE audience_members (
    id INTEGER PRIMARY KEY,
    script_id INTEGER NOT NULL REFERENCES scripts(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    device_token TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL DEFAULT 'waiting' CHECK (status IN ('waiting', 'called', 'writing', 'done', 'skipped')),
    called_at TEXT,
    created_at TEXT NOT NULL,
    UNIQUE (script_id, number)
);

CREATE TABLE writers (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    audience_member_id INTEGER REFERENCES audience_members(id) ON DELETE SET NULL,
    created_at TEXT NOT NULL
);

CREATE TABLE sub_sections (
    id INTEGER PRIMARY KEY,
    section_id INTEGER NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    ordering INTEGER NOT NULL,
    writer_id INTEGER REFERENCES writers(id) ON DELETE SET NULL,
    started_at TEXT,
    ends_at TEXT,
    completed_at TEXT,
    UNIQUE (section_id, ordering)
);

CREATE TABLE lines (
    id INTEGER PRIMARY KEY,
    sub_section_id INTEGER NOT NULL REFERENCES sub_sections(id) ON DELETE CASCADE,
    character_id INTEGER NOT NULL REFERENCES characters(id),
    text TEXT NOT NULL,
    ordering INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE (sub_section_id, ordering)
);

CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
