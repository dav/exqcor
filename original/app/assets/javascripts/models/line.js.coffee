class Exqcor.Models.Line extends Backbone.Model
  methodToURL: =>
    'create': "/plays/#{@get('play_id')}/sections/#{@get('section_id')}/sub_sections/#{@get('sub_section_id')}/lines",
    'delete': "/lines/#{@get('id')}.json"

  sync: (method, model, options) =>
    options = options || {};
    options.url = model.methodToURL()[method.toLowerCase()];
    Backbone.sync(method, model, options)