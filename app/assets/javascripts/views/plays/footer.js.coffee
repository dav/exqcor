class Exqcor.Views.Footer extends Backbone.View
  el: '#footer'
  template: JST['plays/footer']
  initialize: ->
    @collection.bind 'add', @updateCount, @
    @collection.bind 'remove', @updateCount, @
  render: ->
    count = @collection.length
    $(@el).html(@template({playsCount: count}))
    @
  updateCount: ->
    @$('#plays-count').text(@collection.length)