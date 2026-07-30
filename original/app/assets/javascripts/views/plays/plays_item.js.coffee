class Exqcor.Views.PlaysItem extends Backbone.View
  tagName: "tr" 
  template: JST['plays/item']
  
  events:
    'click a.remove-play': 'removePlay'
    'click a.edit-play': 'allowEditPlay'
    
  initialize: ->
    @model.bind 'destroy', @remove, @
    
  render: -> 
    $(@el).html(@template(play: @model))
    @
    
  removePlay: ->
    @model.destroy()
    
  allowEditPlay: ->
    @model.destroy()    