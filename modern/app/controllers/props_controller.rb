class PropsController < ApplicationController
  before_action :set_prop, only: %i[show edit update destroy]

  def index
    props = if params[:section_id]
      Section.find(params[:section_id]).props
    else
      Prop.all
    end
    @props = props
    respond_to do |format|
      format.html
      format.json { render json: @props }
    end
  end

  def show
    respond_to do |format|
      format.html
      format.json { render json: @prop }
    end
  end

  def new
    @prop = Prop.new(section_id: params[:section_id])
  end

  def edit
  end

  def create
    prop = Prop.new(prop_params)
    prop.section_id = params[:section_id] if params[:section_id]

    if prop.save
      respond_to do |format|
        format.html { redirect_to prop_path(prop), notice: "Prop created." }
        format.json { render json: prop, status: :created }
      end
    else
      respond_to do |format|
        format.html do
          @prop = prop
          render :new, status: :unprocessable_entity
        end
        format.json { render_errors(prop) }
      end
    end
  end

  def update
    if @prop.update(prop_params)
      respond_to do |format|
        format.html { redirect_to prop_path(@prop), notice: "Prop updated." }
        format.json { render json: @prop }
      end
    else
      respond_to do |format|
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render_errors(@prop) }
      end
    end
  end

  def destroy
    @prop.destroy
    respond_to do |format|
      format.html { redirect_to props_path, notice: "Prop deleted." }
      format.json { head :no_content }
    end
  end

  private

  def set_prop
    @prop = Prop.find(params[:id])
  end

  def prop_params
    params.require(:prop).permit(:name, :description, :section_id)
  end
end
