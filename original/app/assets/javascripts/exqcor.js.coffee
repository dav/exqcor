window.Exqcor =
	Models: {}
	Collections: {}
	Views: {}
	Routers: {}
	unsavedInput: false
	init: ->
      new Exqcor.Routers.SubSections
      Backbone.history.start()
    handleWriterUnload: ->
      if $('#add-line').val().length > 0
        "You have unsaved changes on this page. If you leave without saving them they will be lost."
      else
        null
      

# window.onbeforeunload = ->
#   Exqcor.handleUnload()


$(document).ready ->
  Exqcor.init()



  # disable globally this for now
  # because need to make sure all pages with a non-ajax save handle
  # it properly (turn off check when user is submitting)
  # $(":input").change ->
  #   console.log 'input changing'
  #   Exqcor.unsavedInput = true # trigers change in all input fields including text type
  #   null
