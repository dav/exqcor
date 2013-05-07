class Exqcor.Views.LinesItem extends Backbone.View
  tagName: "tr" 
  template: JST['lines/item']
  
  events:
    'click a.remove-line': 'removeModel'
    'click a.edit-line': 'allowEditModel'
    
  initialize: ->
    @model.bind 'destroy', @remove, @
    
  render: -> 
    $(@el).html(@template(line: @model))
    @
    
  removeModel: ->
    @model.destroy()
    
  allowEditModel: ->
    @model.destroy()