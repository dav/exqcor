class SubSectionsController < ApplicationController
  # GET /sub_sections/new
  # GET /sub_sections/new.json
  def new
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:section_id])
    index = @section.next_sub_section_index
    @sub_section = @section.sub_sections.create :ordering => index

    respond_to do |format|
      #format.html { redirect_to edit_play_section_sub_section_path(@sub_section), notice: 'Sub section was successfully created.' }
      format.html { redirect_to "/plays/#{@play.id}/sections/#{@section.id}/sub_sections/#{@sub_section.id}/edit", notice: 'Sub section was successfully created.' }
      format.json { render json: @sub_section }
    end
  end

  def show
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:section_id])
    @sub_section = SubSection.find(params[:id])

    respond_to do |format|
      format.html # show.html.erb
      format.json { render json: @sub_section }
    end
  end

  # GET /sub_sections/1/edit
  def edit
    @play = Play.find(params[:play_id])
    @section = Section.find(params[:section_id])
    @sub_section = SubSection.find(params[:id])
  end

  # POST /sub_sections
  # POST /sub_sections.json
  def create
    @play = Play.find(params[:play_id])
    @section = @play.section
    @sub_section = SubSection.new(params[:sub_section])
    

    respond_to do |format|
      if @sub_section.save
        format.html { redirect_to @sub_section, notice: 'Sub section was successfully created.' }
        format.json { render json: @sub_section, status: :created, location: @sub_section }
      else
        format.html { render action: "new" }
        format.json { render json: @sub_section.errors, status: :unprocessable_entity }
      end
    end
  end

  # PUT /sub_sections/1
  # PUT /sub_sections/1.json
  def update
    @sub_section = SubSection.find(params[:id])

    respond_to do |format|
      if @sub_section.update_attributes(params[:sub_section])
        format.html { redirect_to @sub_section, notice: 'Sub section was successfully updated.' }
        format.json { head :no_content }
      else
        format.html { render action: "edit" }
        format.json { render json: @sub_section.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /sub_sections/1
  # DELETE /sub_sections/1.json
  def destroy
    @sub_section = SubSection.find(params[:id])
    @sub_section.destroy

    respond_to do |format|
      format.html { redirect_to sub_sections_url }
      format.json { head :no_content }
    end
  end
end
