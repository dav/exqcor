Rails.application.routes.draw do
  resources :users, only: [:new, :create]
  resource :session
  resources :passwords, param: :token

  namespace :admin do
    resources :users
  end
  concern :admin_setup do
    resources :scripts do
      resources :sections, only: [:index, :new, :create]
      resources :characters, only: [:index, :new, :create]
    end

    resources :sections do
      resources :sub_sections, only: [:index]
      resources :props, only: [:index, :new, :create]
      resources :character_sections, only: [:index, :new, :create]
    end

    resources :sub_sections, only: [:index, :show, :edit, :update] do
      resources :lines, only: [:create]
    end

    resources :characters
    resources :props
    resources :character_sections
  end

  concerns :admin_setup

  scope :api, defaults: { format: :json } do
    post "sign_up", to: "api/auth#sign_up"
    post "sign_in", to: "api/auth#sign_in"
    delete "sign_out", to: "api/auth#sign_out"
    concerns :admin_setup
  end

  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Reveal health status on /up that returns 200 if the app boots with no exceptions, otherwise 500.
  # Can be used by load balancers and uptime monitors to verify that the app is live.
  get "up" => "rails/health#show", as: :rails_health_check

  # Render dynamic PWA files from app/views/pwa/* (remember to link manifest in application.html.erb)
  # get "manifest" => "rails/pwa#manifest", as: :pwa_manifest
  # get "service-worker" => "rails/pwa#service_worker", as: :pwa_service_worker

  # Defines the root path route ("/")
  root "scripts#index"
end
