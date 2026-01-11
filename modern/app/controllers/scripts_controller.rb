class ScriptsController < ApplicationController
  before_action :set_script, only: %i[show edit update destroy]

  def index
    @scripts = Script.all
    respond_to do |format|
      format.html
      format.json { render json: @scripts }
    end
  end

  def show
    respond_to do |format|
      format.html
      format.json { render json: @script }
    end
  end

  def new
    @script = Script.new
  end

  def edit
  end

  def create
    script = Script.new(script_params)
    if script.save
      script.ensure_vosd!
      respond_to do |format|
        format.html { redirect_to script_path(script), notice: "Script created." }
        format.json { render json: script, status: :created }
      end
    else
      respond_to do |format|
        format.html do
          @script = script
          render :new, status: :unprocessable_entity
        end
        format.json { render_errors(script) }
      end
    end
  end

  def update
    if @script.update(script_params)
      respond_to do |format|
        format.html { redirect_to script_path(@script), notice: "Script updated." }
        format.json { render json: @script }
      end
    else
      respond_to do |format|
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render_errors(@script) }
      end
    end
  end

  def destroy
    @script.destroy
    respond_to do |format|
      format.html { redirect_to scripts_path, notice: "Script deleted." }
      format.json { head :no_content }
    end
  end

  private

  def set_script
    @script = Script.find(params[:id])
  end

  def script_params
    params.require(:script).permit(:title, :description)
  end
end
