export type User = {
  id: string;
  email: string;
  display_name: string;
  primary_culture: string;
  created_at: string;
  updated_at: string;
};

export type FamilyMember = {
  id: string;
  user_id: string;
  first_name: string;
  last_name: string;
  display_name: string;
  gender: string;
  birth_date: string | null;
  death_date: string | null;
  birth_place: string;
  biography: string;
  is_living: boolean;
  primary_language: string;
  created_at: string;
  updated_at: string;
};

export type Relationship = {
  id: string;
  user_id: string;
  source_member_id: string;
  target_member_id: string;
  relationship_type: string;
  created_at: string;
};

export type Story = {
  id: string;
  user_id: string;
  family_member_id?: string | null;
  title: string;
  content: string;
  source_type: string;
  source_language: string;
  summary: string;
  occurred_at?: string | null;
  created_at: string;
  updated_at: string;
};

export type WisdomExtract = {
  id: string;
  story_id: string;
  excerpt: string;
  wisdom_type: string;
  language: string;
  meaning: string;
  confidence: number;
  created_at: string;
};
