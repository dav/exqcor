# encoding: UTF-8
# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended to check this file into your version control system.

ActiveRecord::Schema.define(:version => 20130917141440) do

  create_table "characters", :force => true do |t|
    t.string   "name"
    t.string   "description"
    t.integer  "play_id"
    t.datetime "created_at",  :null => false
    t.datetime "updated_at",  :null => false
  end

  create_table "characters_sections", :id => false, :force => true do |t|
    t.integer "character_id"
    t.integer "section_id"
  end

  add_index "characters_sections", ["character_id", "section_id"], :name => "index_characters_sections_on_character_id_and_section_id"
  add_index "characters_sections", ["section_id"], :name => "index_characters_sections_on_section_id"

  create_table "lines", :force => true do |t|
    t.text     "text"
    t.integer  "ordering"
    t.integer  "sub_section_id"
    t.integer  "character_id"
    t.datetime "created_at",     :null => false
    t.datetime "updated_at",     :null => false
  end

  create_table "plays", :force => true do |t|
    t.string   "title"
    t.string   "description"
    t.datetime "created_at",  :null => false
    t.datetime "updated_at",  :null => false
  end

  create_table "props", :force => true do |t|
    t.string   "name"
    t.string   "description"
    t.integer  "section_id"
    t.datetime "created_at",  :null => false
    t.datetime "updated_at",  :null => false
  end

  create_table "sections", :force => true do |t|
    t.string   "name"
    t.string   "description"
    t.integer  "ordering"
    t.integer  "play_id"
    t.datetime "created_at",  :null => false
    t.datetime "updated_at",  :null => false
  end

  create_table "sub_sections", :force => true do |t|
    t.integer  "ordering"
    t.integer  "section_id"
    t.datetime "created_at", :null => false
    t.datetime "updated_at", :null => false
  end

end
