class Exqcor.Views.SubSectionsWriter extends Backbone.View
  el: '#sub_section_writer'
  template: JST['sub_sections/writer']

  events: 
    'keypress #add-line': 'createOnEnter'
    
  initialize: ->
    console.log 'in sub_sections writer init with model of type '+@model.constructor.name
    console.log "this is a #{@.constructor.name}"

  render: ->
    console.log 'in render'
    $(@el).html(@template())
    
  createOnEnter: (event) ->
    return if event.keyCode != 13
    @collection.create title: @$('#add-play').val()
    @$('#add-play').val('')
  
  addLine: (line) ->
    view = new Exqcor.Views.LinesItem model: line
    @$('#lines').append(view.render().el)
    @