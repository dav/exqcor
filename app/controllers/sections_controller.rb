class SectionsController < ApplicationController
  # GET /Sections/new
  # GET /Sections/new.json
  def new
    @play = Play.find(params[:play_id])
    @section = @play.sections.build

    respond_to do |format|
      format.html # new.html.erb
      format.json { render json: @section }
    end
  end

  # GET /Sections/1/edit
  def edit
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:id])
  end

  def show
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:id])
    if params[:character_id]
      @character = Character.find(params[:character_id])
    end
  end

  def script_package
    @section = Section.find(params[:section_id])

    respond_to do |format|
      format.html
    end
  end
  
  def script_for_character
    @section = Section.find(params[:section_id])
    @character = Character.find(params[:character_id])

    respond_to do |format|
      format.html
    end
  end

  def script_for_actor
    @section = Section.find(params[:section_id])
    @actor = Actor.find(params[:actor_id])

    respond_to do |format|
      format.html
    end
  end

  # POST /Sections
  # POST /Sections.json
  def create
    @play = Play.find(params[:play_id])
    @section = Section.new(params[:section])
    @section.play = @play

    respond_to do |format|
      if @section.save
        format.html { redirect_to edit_play_url(@play), notice: 'section was successfully created.' }
        format.json { render json: @section, status: :created, location: @section }
      else
        format.html { render action: "new" }
        format.json { render json: @section.errors, status: :unprocessable_entity }
      end
    end
  end

  # PUT /Sections/1
  # PUT /Sections/1.json
  def update
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:id])

    respond_to do |format|
      if @section.update_attributes(params[:section])
        format.html { redirect_to edit_play_url(@play), notice: 'section was successfully updated.' }
        format.json { head :no_content }
      else
        format.html { redirect_to edit_play_url(@play) }
        format.json { render json: @section.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /Sections/1
  # DELETE /Sections/1.json
  def destroy
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:id])
    @section.destroy

    respond_to do |format|
      format.html { redirect_to edit_play_url(@play) }
      format.json { head :no_content }
    end
  end
  
  def characters
    @section = Section.find(params[:section_id])
    @character = Character.find(params[:character_id])
    character_param = 'character_'+params[:character_id]
    selected = params[character_param] == 'selected'
    on_stage_param = params[:on_stage]
    if @section && @character
      if params[character_param]
        if selected
          logger.debug "ADD #{@character.name}"
          @section.characters << @character
        else
          logger.debug "REMOVE #{@character.name}"
          @section.characters.delete(@character)
        end
      end

      if on_stage_param
        cs = CharacterSection.where(character_id: @character.id, section_id: @section.id).first
        cs.update_attribute(:on_stage, (on_stage_param=='true')) unless cs.nil?
      end
    end
    
    @section.reload
    respond_to do |format|
      format.html { render json: @section.characters.map(&:id) }
    end
  end
end
