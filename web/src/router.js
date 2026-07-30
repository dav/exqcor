import Home from './routes/Home.svelte'
import Dashboard from './routes/admin/Dashboard.svelte'
import Scripts from './routes/admin/Scripts.svelte'
import ScriptEditor from './routes/admin/ScriptEditor.svelte'
import RunOfShow from './routes/admin/RunOfShow.svelte'
import QRPage from './routes/admin/QRPage.svelte'
import WritingStation from './routes/write/WritingStation.svelte'
import Queue from './routes/audience/Queue.svelte'
import Program from './routes/audience/Program.svelte'
import ScriptView from './routes/script/ScriptView.svelte'
import NotFound from './routes/NotFound.svelte'

export const routes = {
  '/': Home,
  '/admin': Dashboard,
  '/admin/scripts': Scripts,
  '/admin/scripts/:id': ScriptEditor,
  '/admin/show': RunOfShow,
  '/admin/qr': QRPage,
  '/write': WritingStation,
  '/audience': Queue,
  '/program': Program,
  '/script/:id': ScriptView,
  '/script/:id/actor/:actorId': ScriptView,
  '/script/:id/character/:characterId': ScriptView,
  '*': NotFound,
}
