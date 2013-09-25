class Exqcor.Views.CharacterOptionItem extends Backbone.View
  tagName: "option"
  
  attributes: ->
    value: @model.get('id')
  
  template: JST['characters/option']
  
  render: -> 
    if @model.get('name') == 'VOSD'
      @model.set('description', 'Voice of Stage Directions') 
    $(@el).html(@template(character: @model))
    @
