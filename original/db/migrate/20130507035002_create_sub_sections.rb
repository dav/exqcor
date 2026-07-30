class CreateSubSections < ActiveRecord::Migration
  def change
    create_table :sub_sections do |t|
      t.integer :ordering
      t.integer :section_id

      t.timestamps
    end
  end
end
