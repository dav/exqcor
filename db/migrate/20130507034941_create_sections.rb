class CreateSections < ActiveRecord::Migration
  def change
    create_table :sections do |t|
      t.string :name
      t.string :description
      t.integer :ordering
      t.string :play_id

      t.timestamps
    end
  end
end
