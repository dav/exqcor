class Exqcor.Views.LinesItem extends Backbone.View
  tagName: "tr" 
  template: JST['lines/item']
  
  events:
    'click input.remove-line': 'removeModel'
    #'click a.edit-line': 'allowEditModel'
    
  initialize: ->
    @model.bind 'destroy', @remove, @
    
  render: ->
    $(@el).html(@template(line: @model, character: @options['character']))
    @
    
  removeModel: ->
    $('#line-'+@model.get('id')).css({ 'background-color': 'rgba(255,0,0,.5)'})
    if confirm("Are you sure you want to delete this line?")
      @model.destroy()
    $('#line-'+@model.get('id')).css({ 'background-color': 'white'})
    
  #allowEditModel: ->
  #  @model.destroy()