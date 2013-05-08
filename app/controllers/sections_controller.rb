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
    @section = section.find(params[:id])
  end

  def show
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:id])
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
end
