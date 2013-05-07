window.Exqcor =
	Models: {}
	Collections: {}
	Views: {}
	Routers: {}
	init: ->
	  #new Exqcor.Routers.Plays
	  new Exqcor.Routers.SubSections
	  Backbone.history.start()
	
$(document).ready ->
	Exqcor.init()