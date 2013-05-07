class Exqcor.Views.SubSectionsWriter extends Backbone.View
  el: '#sub_section_writer'
  
  template: JST['sub_sections/writer']

  events: 
    'keypress #add-line': 'createOnEnter'
    
  initialize: ->
    console.log "this is a #{@.constructor.name}"
    console.log 'in sub_sections writer init with model of type '+@model.constructor.name
    console.log 'in sub_sections writer init with collection of type '+@collection.constructor.name
    console.log 'in sub_sections writer init with collection of size '+@collection.size()
    console.log 'in sub_sections writer init with play_id of '+@model.get('play_id')
    @collection.bind 'reset', @render, @
    @collection.bind 'add', @addLine, @

  render: ->
    console.log 'in render'
    $(@el).html(@template())
    @collection.each (line) =>
      console.log "LINE #{line.get('character_id')} #{line.get('text')}"
      view = new Exqcor.Views.LinesItem model: line
      @$('#lines').append(view.render().el)
    @
    
  createOnEnter: (event) ->
    return if event.keyCode != 13
    
    @collection.create text: @$('#add-line').val(), \
      character_id: @$('#add-line-character_id').val(), \
      play_id: @model.get('play_id'), \
      section_id: @model.get('section_id'), 
      sub_section_id: @model.get('sub_section_id')
    @$('#add-line').val('')
  
  addLine: (line) ->
    view = new Exqcor.Views.LinesItem model: line
    @$('#lines').append(view.render().el)
    @