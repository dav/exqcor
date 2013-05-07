class CreateProps < ActiveRecord::Migration
  def change
    create_table :props do |t|
      t.string :name
      t.string :description
      t.string :section_id

      t.timestamps
    end
  end
end
