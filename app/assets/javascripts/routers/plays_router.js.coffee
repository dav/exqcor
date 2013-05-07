class Exqcor.Routers.Plays extends Backbone.Router
  routes:
    '': 'index'
    
  fetch_success: (collection, response) =>
    #console.log "we fetched a #{collection.constructor.name} with response #{JSON.stringify(collection.toJSON())}"
    collection.reset response
    
  index: ->
    plays = new Exqcor.Collections.Plays
    new Exqcor.Views.PlaysIndex collection: plays # this should attach the view to the collection
    plays.fetch
      success: @fetch_success
    # console.log 'end of index and plays collection looks like: (next line)'
    # console.log plays

