class Exqcor.Views.CharacterOptionItem extends Backbone.View
  tagName: "option"
  
  attributes: ->
    value: @model.get('id')
  
  template: JST['characters/option']
  
  render: -> 
    $(@el).html(@template(character: @model))
    @
