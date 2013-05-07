class PropsController < ApplicationController
  # GET /props/new
  # GET /props/new.json
  def new
    @section = Section.find(params[:section_id])
    @play = @section.play
    @prop = @section.props.build

    respond_to do |format|
      format.html # new.html.erb
      format.json { render json: @prop }
    end
  end

  # GET /props/1/edit
  def edit
    @prop = Prop.find(params[:id])
  end

  # POST /props
  # POST /props.json
  def create
    @section = Section.find(params[:section_id])
    @play = @section.play
    @prop = Prop.new(params[:prop])
    @prop.section = @section

    respond_to do |format|
      if @prop.save
        format.html { redirect_to edit_play_url(@play), notice: 'Prop was successfully created.' }
        format.json { render json: @prop, status: :created, location: @prop }
      else
        format.html { render action: "new" }
        format.json { render json: @prop.errors, status: :unprocessable_entity }
      end
    end
  end

  # PUT /props/1
  # PUT /props/1.json
  def update
    @prop = Prop.find(params[:id])

    respond_to do |format|
      if @prop.update_attributes(params[:prop])
        format.html { redirect_to edit_play_url(@play), notice: 'Prop was successfully updated.' }
        format.json { head :no_content }
      else
        format.html { redirect_to edit_play_url(@play), }
        format.json { render json: @prop.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /props/1
  # DELETE /props/1.json
  def destroy
    @play = Play.find(params[:play_id])
    @prop = Prop.find(params[:id])
    @prop.destroy

    respond_to do |format|
      format.html { redirect_to edit_play_url(@play) }
      format.json { head :no_content }
    end
  end
end
