class CharactersController < ApplicationController
  before_action :set_character, only: %i[show edit update destroy]

  def index
    characters = if params[:script_id]
      Script.find(params[:script_id]).characters
    else
      Character.all
    end
    @characters = characters
    respond_to do |format|
      format.html
      format.json { render json: @characters }
    end
  end

  def show
    respond_to do |format|
      format.html
      format.json { render json: @character }
    end
  end

  def new
    @character = Character.new(script_id: params[:script_id])
  end

  def edit
  end

  def create
    character = Character.new(character_params)
    character.script_id = params[:script_id] if params[:script_id]

    if character.save
      respond_to do |format|
        format.html { redirect_to character_path(character), notice: "Character created." }
        format.json { render json: character, status: :created }
      end
    else
      respond_to do |format|
        format.html do
          @character = character
          render :new, status: :unprocessable_entity
        end
        format.json { render_errors(character) }
      end
    end
  end

  def update
    if @character.update(character_params)
      respond_to do |format|
        format.html { redirect_to character_path(@character), notice: "Character updated." }
        format.json { render json: @character }
      end
    else
      respond_to do |format|
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render_errors(@character) }
      end
    end
  end

  def destroy
    @character.destroy
    respond_to do |format|
      format.html { redirect_to characters_path, notice: "Character deleted." }
      format.json { head :no_content }
    end
  end

  private

  def set_character
    @character = Character.find(params[:id])
  end

  def character_params
    params.require(:character).permit(:name, :description, :script_id, :actor_id, :role)
  end
end
