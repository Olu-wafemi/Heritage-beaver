CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    primary_culture TEXT NOT NULL DEFAULT 'yoruba',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS family_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    first_name TEXT NOT NULL,
    last_name TEXT,
    display_name TEXT NOT NULL,
    gender TEXT,
    birth_date DATE,
    death_date DATE,
    birth_place TEXT,
    biography TEXT,
    is_living BOOLEAN NOT NULL DEFAULT TRUE,
    primary_language TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_member_id UUID NOT NULL REFERENCES family_members(id) ON DELETE CASCADE,
    target_member_id UUID NOT NULL REFERENCES family_members(id) ON DELETE CASCADE,
    relationship_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT relationships_no_self_link CHECK (source_member_id <> target_member_id)
);

CREATE TABLE IF NOT EXISTS stories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_member_id UUID REFERENCES family_members(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    source_type TEXT NOT NULL,
    source_language TEXT NOT NULL DEFAULT 'en',
    summary TEXT,
    occurred_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wisdom_extracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    excerpt TEXT NOT NULL,
    wisdom_type TEXT NOT NULL,
    language TEXT NOT NULL DEFAULT 'en',
    meaning TEXT,
    confidence DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS myth_chapters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    theme TEXT NOT NULL,
    chapter_type TEXT NOT NULL,
    narrative TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ancestor_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_member_id UUID REFERENCES family_members(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    persona_summary TEXT NOT NULL,
    voice_style TEXT,
    guidance_style TEXT,
    cultural_context TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS guidance_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ancestor_profile_id UUID REFERENCES ancestor_profiles(id) ON DELETE SET NULL,
    myth_chapter_id UUID REFERENCES myth_chapters(id) ON DELETE SET NULL,
    prompt TEXT NOT NULL,
    response TEXT NOT NULL,
    session_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS media_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_type TEXT NOT NULL,
    storage_key TEXT NOT NULL,
    source_kind TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS privacy_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    allow_family_sharing BOOLEAN NOT NULL DEFAULT FALSE,
    allow_voice_cloning BOOLEAN NOT NULL DEFAULT FALSE,
    allow_ancestor_simulation BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_family_members_user_id ON family_members (user_id);
CREATE INDEX IF NOT EXISTS idx_relationships_user_id ON relationships (user_id);
CREATE INDEX IF NOT EXISTS idx_relationships_source_target ON relationships (source_member_id, target_member_id);
CREATE INDEX IF NOT EXISTS idx_stories_user_id ON stories (user_id);
CREATE INDEX IF NOT EXISTS idx_stories_family_member_id ON stories (family_member_id);
CREATE INDEX IF NOT EXISTS idx_wisdom_extracts_story_id ON wisdom_extracts (story_id);
CREATE INDEX IF NOT EXISTS idx_myth_chapters_user_id ON myth_chapters (user_id);
CREATE INDEX IF NOT EXISTS idx_ancestor_profiles_user_id ON ancestor_profiles (user_id);
CREATE INDEX IF NOT EXISTS idx_guidance_sessions_user_id ON guidance_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_media_assets_user_id ON media_assets (user_id);
