class SubSectionsController < ApplicationController
  before_action :set_sub_section, only: %i[show edit update]

  def index
    sub_sections = if params[:section_id]
      Section.find(params[:section_id]).sub_sections
    else
      SubSection.all
    end
    @sub_sections = sub_sections
    respond_to do |format|
      format.html
      format.json { render json: @sub_sections }
    end
  end

  def show
    respond_to do |format|
      format.html
      format.json { render json: @sub_section }
    end
  end

  def edit
    @section = @sub_section.section
    @script = @section.script
    @line = Line.new
    @last_sub_section = @section.sub_sections.where("ordering < ?", @sub_section.ordering).order(ordering: :desc).first
    @last_line = @last_sub_section&.lines&.order(ordering: :desc)&.first
    vosd = @script.VOSD
    @characters = [vosd, *@section.characters.where.not(id: vosd&.id).order(:name)].compact
    @writers = Writer.order(:name)
  end

  def update
    if @sub_section.update(sub_section_params)
      respond_to do |format|
        format.html { redirect_to edit_sub_section_path(@sub_section), notice: "Writer updated." }
        format.json { render json: @sub_section }
      end
    else
      respond_to do |format|
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render_errors(@sub_section) }
      end
    end
  end

  private

  def set_sub_section
    @sub_section = SubSection.find(params[:id])
  end

  def sub_section_params
    params.require(:sub_section).permit(:writer_id)
  end
end
