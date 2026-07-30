class Exqcor.Models.SubSection extends Backbone.Model
	url: =>
		"/plays/#{@get('play_id')}/sections/#{@get('section_id')}/sub_sections/#{@get('sub_section_id')}"
