class CreateScripts < ActiveRecord::Migration[8.0]
  def change
    create_table :scripts do |t|
      t.string :title
      t.string :description

      t.timestamps
    end
  end
end
