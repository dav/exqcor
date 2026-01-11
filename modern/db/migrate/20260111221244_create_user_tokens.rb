class CreateUserTokens < ActiveRecord::Migration[8.0]
  def change
    create_table :user_tokens do |t|
      t.references :user, null: false, foreign_key: true
      t.string :token_digest, null: false
      t.string :name
      t.datetime :last_used_at
      t.datetime :expires_at

      t.timestamps
    end

    add_index :user_tokens, :token_digest, unique: true
  end
end
