class Exqcor.Views.PlaysIndex extends Backbone.View
  el: '#app'
  
  template: JST['plays/index']
  
  events: 
    'keypress #add-play': 'createOnEnter'
    
  initialize: ->
    console.log 'in plays index init with collection of type '+@collection.constructor.name
    console.log "this is a #{@.constructor.name}"
    @collection.bind 'reset', @render, @
    @collection.bind 'add', @addPlay, @
    
  render: ->
    console.log 'in render'
    $(@el).html(@template())
    
    footerView = new Exqcor.Views.Footer collection: @collection
    footerView.render()
    
    @collection.each (play) =>
      console.log "PLAY #{play.get('title')}"
      view = new Exqcor.Views.PlaysItem model: play
      @$('#plays').append(view.render().el)
    @
    
  createOnEnter: (event) ->
    return if event.keyCode != 13
    @collection.create title: @$('#add-play').val()
    @$('#add-play').val('')
  
  addPlay: (play) ->
    view = new Exqcor.Views.PlaysItem model: play
    @$('#plays').append(view.render().el)
    @