window.Exqcor =
	Models: {}
	Collections: {}
	Views: {}
	Routers: {}
	init: ->
	  new Exqcor.Routers.Plays
	  Backbone.history.start()
	
$(document).ready ->
	Exqcor.init()