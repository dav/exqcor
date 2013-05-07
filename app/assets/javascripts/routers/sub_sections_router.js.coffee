class Exqcor.Routers.SubSections extends Backbone.Router
	routes:
	  'write/play/:play_id/section/:section_id/sub_section/:sub_section_id': 'writeSubSection'
  
	fetch_success: (collection, response) =>
	  console.log "we fetched a #{collection.constructor.name} with response #{JSON.stringify(collection.toJSON())}"
	  collection.reset response

	fetch_sub_section_success: (model, response) =>
	  console.log "we fetched a #{model.constructor.name}"
	  console.log "with response #{JSON.stringify(model.toJSON())}"
	  #lines = new Exqcor.Collections.Lines
	  view = new Exqcor.Views.SubSectionsWriter model: model
	  view.render()
	  # subsections.fetch
	  #   success: @fetch_success
  
	writeSubSection: (play_id, section_id, sub_section_id) ->
	  console.log "START #{play_id}, #{section_id}, #{sub_section_id}, "
	  subsection = new Exqcor.Models.SubSection play_id: play_id, section_id: section_id, sub_section_id: sub_section_id
	  console.log 'subsection url is '+subsection.url()
	  subsection.fetch
	    success: @fetch_sub_section_success

	  console.log 'end of index and subsections collection looks like: (next line)'
	  console.log subsection
