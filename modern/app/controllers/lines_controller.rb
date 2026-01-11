class LinesController < ApplicationController
  def create
    sub_section = SubSection.find(params[:sub_section_id])
    line = sub_section.lines.new(line_params)
    line.ordering ||= sub_section.next_ordering_index

    if line.save
      respond_to do |format|
        format.html { redirect_to edit_sub_section_path(sub_section), notice: "Line added." }
        format.json { render json: line, status: :created }
      end
    else
      respond_to do |format|
        format.html do
          @sub_section = sub_section
          @section = sub_section.section
          @script = @section.script
          @line = line
          @last_sub_section = @section.sub_sections.where("ordering < ?", @sub_section.ordering).order(ordering: :desc).first
          @last_line = @last_sub_section&.lines&.order(ordering: :desc)&.first
          vosd = @script.VOSD
          @characters = [vosd, *@section.characters.where.not(id: vosd&.id).order(:name)].compact
          @writers = Writer.order(:name)
          render "sub_sections/edit", status: :unprocessable_entity
        end
        format.json { render_errors(line) }
      end
    end
  end

  private

  def line_params
    params.require(:line).permit(:text, :ordering, :character_id)
  end
end
