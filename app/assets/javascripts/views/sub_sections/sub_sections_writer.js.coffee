class Exqcor.Views.SubSectionsWriter extends Backbone.View
  el: '#sub_section_writer'
  
  template: JST['sub_sections/writer']

  events: 
    'click #add-line-button': 'createLine'
    
  initialize: ->
    console.log "this is a #{@.constructor.name}"
    console.log 'in sub_sections writer init with model of type '+@model.constructor.name
    console.log 'in sub_sections writer init with collection of type '+@collection.constructor.name
    console.log 'in sub_sections writer init with collection of size '+@collection.size()
    console.log 'in sub_sections writer init with play_id of '+@model.get('play_id')
    @collection.bind 'reset', @render, @
    @collection.bind 'add', @addLine, @

  render: ->
    $(@el).html(@template())
    # add character choices
    characters = @model.get('characters')
    characters.each (character) =>
      console.log "CHARACTER #{character.get('id')} #{character.get('name')}"
      view = new Exqcor.Views.CharacterOptionItem model: character
      @$('#character-id-select').append(view.render().el)
    # add lines
    @collection.each (line) =>
      character = characters.get(line.get('character_id'))
      console.log "LINE #{character.get('name')} says: #{line.get('text')}"
      view = new Exqcor.Views.LinesItem model: line, character: character
      @$('#lines').append(view.render().el)
    @

  createLine: (event) ->
    console.log('createLine');
    @collection.create \
      text: @$('#add-line').val(), \
      character_id: @$('#character-id-select').val(), \
      play_id: @model.get('play_id'), \
      section_id: @model.get('section_id'), \
      sub_section_id: @model.get('sub_section_id')
    @$('#add-line').val('')
    @$('#add-line-character_id').focus();
    @
  
  addLine: (line) ->
    characters = @model.get('characters')
    character = characters.get(line.get('character_id'))
    view = new Exqcor.Views.LinesItem model: line, character: character
    @$('#lines').append(view.render().el)
    @