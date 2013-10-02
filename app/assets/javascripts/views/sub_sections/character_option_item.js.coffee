class Exqcor.Views.CharacterOptionItem extends Backbone.View
  tagName: "option"
  
  attributes: ->
    value: @model.get('id')
    # the title is the tooltip
    title: @model.get('name')+': '+@model.get('description')
  
  template: JST['characters/option']
  
  render: ->
    if @model.get('name') == 'VOSD'
      @model.set('description', 'Voice of Stage Directions') 
    $(@el).html(@template(character: @model))
    @
