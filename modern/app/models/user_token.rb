class UserToken < ApplicationRecord
  belongs_to :user

  before_validation :set_token_digest, on: :create

  def self.issue!(user, name: nil, expires_at: nil)
    raw_token = generate_token
    create!(
      user: user,
      token_digest: digest(raw_token),
      name: name,
      last_used_at: Time.current,
      expires_at: expires_at || 30.days.from_now
    )
    raw_token
  end

  def self.find_by_token(token)
    return nil if token.blank?

    find_by(token_digest: digest(token))
  end

  def expired?
    expires_at.present? && Time.current > expires_at
  end

  def touch_last_used!
    touch(:last_used_at)
  end

  def self.generate_token
    SecureRandom.base58(48)
  end

  def self.digest(token)
    OpenSSL::Digest::SHA256.hexdigest(token)
  end

  def rotate!
    raw_token = self.class.generate_token
    new_token = self.class.create!(
      user: user,
      token_digest: self.class.digest(raw_token),
      name: name,
      last_used_at: Time.current,
      expires_at: 30.days.from_now
    )
    destroy!
    [raw_token, new_token]
  end

  private

  def set_token_digest
    self.token_digest ||= self.class.digest(self.class.generate_token)
  end
end
