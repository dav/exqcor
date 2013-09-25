# http://localhost:3000/plays/1/sections/4/sub_sections/11/edit#write/play/1/section/4/sub_section/11
class Exqcor.Routers.SubSections extends Backbone.Router
	routes:
	  'write/play/:play_id/section/:section_id/sub_section/:sub_section_id': 'writeSubSection'
  
	#fetch_success: (collection, response) =>
	#  console.log "we fetched a #{collection.constructor.name} with response #{JSON.stringify(collection.toJSON())}"
	#  collection.reset response

	fetch_sub_section_success: (model, response) =>
	  #console.log "we fetched a #{model.constructor.name}"
	  #console.log "with response #{response}"
	  #console.log "with response #{JSON.stringify(model.toJSON())}"
	  lines = new Exqcor.Collections.Lines model.get('lines')
	  netCharactersArray = response['section']['sorted_characters']
	  #console.log "play #{netPlay}"
	  characters = new Exqcor.Collections.Characters netCharactersArray
	  model.set('characters', characters)
	  view = new Exqcor.Views.SubSectionsWriter model: model, collection: lines
	  view.render()
	  # subsections.fetch
	  #   success: @fetch_success
  
	writeSubSection: (play_id, section_id, sub_section_id) ->
	  subsection = new Exqcor.Models.SubSection play_id: play_id, section_id: section_id, sub_section_id: sub_section_id
	  subsection.fetch
	    success: @fetch_sub_section_success
