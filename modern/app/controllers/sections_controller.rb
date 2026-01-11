class SectionsController < ApplicationController
  before_action :set_section, only: %i[show edit update destroy]

  def index
    sections = if params[:script_id]
      Script.find(params[:script_id]).sections
    else
      Section.all
    end
    @sections = sections
    respond_to do |format|
      format.html
      format.json { render json: @sections }
    end
  end

  def show
    respond_to do |format|
      format.html
      format.json { render json: @section }
    end
  end

  def new
    @section = Section.new(script_id: params[:script_id])
  end

  def edit
  end

  def create
    section = Section.new(section_params)
    section.script_id = params[:script_id] if params[:script_id]

    if section.save
      section.ensure_vosd!
      respond_to do |format|
        format.html { redirect_to section_path(section), notice: "Section created." }
        format.json { render json: section, status: :created }
      end
    else
      respond_to do |format|
        format.html do
          @section = section
          render :new, status: :unprocessable_entity
        end
        format.json { render_errors(section) }
      end
    end
  end

  def update
    if @section.update(section_params)
      respond_to do |format|
        format.html { redirect_to section_path(@section), notice: "Section updated." }
        format.json { render json: @section }
      end
    else
      respond_to do |format|
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render_errors(@section) }
      end
    end
  end

  def destroy
    @section.destroy
    respond_to do |format|
      format.html { redirect_to sections_path, notice: "Section deleted." }
      format.json { head :no_content }
    end
  end

  private

  def set_section
    @section = Section.find(params[:id])
  end

  def section_params
    params.require(:section).permit(:name, :description, :ordering, :script_id)
  end
end
