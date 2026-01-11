class CharacterSectionsController < ApplicationController
  before_action :set_character_section, only: %i[show edit update destroy]

  def index
    character_sections = CharacterSection.all
    character_sections = character_sections.where(section_id: params[:section_id]) if params[:section_id]
    character_sections = character_sections.where(character_id: params[:character_id]) if params[:character_id]
    @character_sections = character_sections
    respond_to do |format|
      format.html
      format.json { render json: @character_sections }
    end
  end

  def show
    respond_to do |format|
      format.html
      format.json { render json: @character_section }
    end
  end

  def new
    @character_section = CharacterSection.new(section_id: params[:section_id], character_id: params[:character_id])
  end

  def edit
  end

  def create
    character_section = CharacterSection.new(character_section_params)
    character_section.section_id = params[:section_id] if params[:section_id]

    if character_section.save
      respond_to do |format|
        format.html { redirect_to character_section_path(character_section), notice: "Character section created." }
        format.json { render json: character_section, status: :created }
      end
    else
      respond_to do |format|
        format.html do
          @character_section = character_section
          render :new, status: :unprocessable_entity
        end
        format.json { render_errors(character_section) }
      end
    end
  end

  def update
    if @character_section.update(character_section_params)
      respond_to do |format|
        format.html { redirect_to character_section_path(@character_section), notice: "Character section updated." }
        format.json { render json: @character_section }
      end
    else
      respond_to do |format|
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render_errors(@character_section) }
      end
    end
  end

  def destroy
    @character_section.destroy
    respond_to do |format|
      format.html { redirect_to character_sections_path, notice: "Character section deleted." }
      format.json { head :no_content }
    end
  end

  private

  def set_character_section
    @character_section = CharacterSection.find(params[:id])
  end

  def character_section_params
    params.require(:character_section).permit(:character_id, :section_id, :on_stage)
  end
end
