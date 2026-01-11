# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[8.0].define(version: 2026_01_11_223600) do
  create_table "actors", force: :cascade do |t|
    t.string "name"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "character_sections", force: :cascade do |t|
    t.integer "character_id", null: false
    t.integer "section_id", null: false
    t.boolean "on_stage"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["character_id", "section_id"], name: "index_character_sections_on_character_id_and_section_id", unique: true
    t.index ["character_id"], name: "index_character_sections_on_character_id"
    t.index ["section_id"], name: "index_character_sections_on_section_id"
  end

  create_table "characters", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.integer "script_id", null: false
    t.integer "actor_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.string "role", default: "character", null: false
    t.index ["actor_id"], name: "index_characters_on_actor_id"
    t.index ["script_id", "name"], name: "index_characters_on_script_id_and_name", unique: true
    t.index ["script_id"], name: "index_characters_on_script_id"
    t.index ["script_id"], name: "index_characters_on_script_id_unique_vosd", unique: true, where: "role = 'vosd'"
  end

  create_table "lines", force: :cascade do |t|
    t.text "text"
    t.integer "ordering"
    t.integer "sub_section_id", null: false
    t.integer "character_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["character_id"], name: "index_lines_on_character_id"
    t.index ["sub_section_id"], name: "index_lines_on_sub_section_id"
  end

  create_table "props", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.integer "section_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["section_id"], name: "index_props_on_section_id"
  end

  create_table "scripts", force: :cascade do |t|
    t.string "title"
    t.string "description"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  create_table "sections", force: :cascade do |t|
    t.string "name"
    t.string "description"
    t.integer "ordering"
    t.integer "script_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["script_id"], name: "index_sections_on_script_id"
  end

  create_table "sessions", force: :cascade do |t|
    t.integer "user_id", null: false
    t.string "ip_address"
    t.string "user_agent"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["user_id"], name: "index_sessions_on_user_id"
  end

  create_table "sub_sections", force: :cascade do |t|
    t.integer "ordering"
    t.integer "section_id", null: false
    t.integer "writer_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["section_id"], name: "index_sub_sections_on_section_id"
    t.index ["writer_id"], name: "index_sub_sections_on_writer_id"
  end

  create_table "user_tokens", force: :cascade do |t|
    t.integer "user_id", null: false
    t.string "token_digest", null: false
    t.string "name"
    t.datetime "last_used_at"
    t.datetime "expires_at"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["token_digest"], name: "index_user_tokens_on_token_digest", unique: true
    t.index ["user_id"], name: "index_user_tokens_on_user_id"
  end

  create_table "users", force: :cascade do |t|
    t.string "email_address", null: false
    t.string "password_digest", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.boolean "admin", default: false, null: false
    t.index ["email_address"], name: "index_users_on_email_address", unique: true
  end

  create_table "writers", force: :cascade do |t|
    t.string "name"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
  end

  add_foreign_key "character_sections", "characters"
  add_foreign_key "character_sections", "sections"
  add_foreign_key "characters", "actors"
  add_foreign_key "characters", "scripts"
  add_foreign_key "lines", "characters"
  add_foreign_key "lines", "sub_sections"
  add_foreign_key "props", "sections"
  add_foreign_key "sections", "scripts"
  add_foreign_key "sessions", "users"
  add_foreign_key "sub_sections", "sections"
  add_foreign_key "sub_sections", "writers"
  add_foreign_key "user_tokens", "users"
end
