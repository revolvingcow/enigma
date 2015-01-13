package main

const (
	templateIndex string = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Enigma Login</title>

    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css" />
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap-theme.min.css" />
    <style type="text/css">
        body {
            padding-top: 40px;
            padding-bottom: 40px;
            background-color: #eee;
        }

        .form-signin {
            max-width: 330px;
            padding: 15px;
            margin: 0 auto;
        }
        .form-signin .form-signin-heading,
        .form-signin .checkbox {
            margin-bottom: 10px;
        }
        .form-signin .checkbox {
            font-weight: normal;
        }
        .form-signin .form-control {
            position: relative;
            height: auto;
            -webkit-box-sizing: border-box;
            -moz-box-sizing: border-box;
            box-sizing: border-box;
            padding: 10px;
            font-size: 16px;
        }
        .form-signin .form-control:focus {
            z-index: 2;
        }
        .form-signin input[type="email"] {
            margin-bottom: -1px;
            border-bottom-right-radius: 0;
            border-bottom-left-radius: 0;
        }
        .form-signin input[type="password"] {
            margin-bottom: 10px;
            border-top-left-radius: 0;
            border-top-right-radius: 0;
        }
    </style>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
        <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>
    <div class="container">
        <form action="/" method="post" class="form-signin">
            <h2 class="form-signin-heading">Please sign in</h2>
            <label for="profile" class="sr-only">Email address</label>
            <input type="text" id="profile" name="profile" class="form-control" placeholder="Username" required autofocus>
            <label for="p" class="sr-only">Password</label>
            <input type="password" id="p" name="p" class="form-control" placeholder="Password" required>
            <div class="checkbox">
                <label>
                    <input type="checkbox" value="remember-me"> Remember me
                </label>
            </div>
            <button class="btn btn-lg btn-primary btn-block" type="submit">Sign in</button>
        </form>
    </div>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/js/bootstrap.min.js"></script>
</body>
</html>`

	templateBook string = `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Enigma</title>

    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css" />
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap-theme.min.css" />
    <style type="text/css">
        body {
            padding-top: 40px;
            padding-bottom: 40px;
            background-color: #fff;
        }
        h2 {
            margin-bottom: 1em;
        }
        td {
            text-align: left;
            vertical-align: middle !important;
        }
        .tab-content > .tab-pane {
            padding: 1em;
        }
    </style>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
        <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>
<body>
    <div class="container">
        <h1>Enigma</h1>
        <h2><small>Your personal password safe and generator</small></h2>

		{{ $profile := .Profile }}
		{{ $passphrase := .Passphrase }}

        <div role="tabpanel">
            <ul class="nav nav-tabs" role="tablist">
                <li role="presentation" class="active"><a href="#passwords" aria-controls="passwords" role="tab" data-toggle="tab">Passwords</a></li>
                <li role="presentation"><a href="#add" aria-controls="Add" role="tab" data-toggle="tab">Add</a></li>
                <li role="presentation"><a href="#settings" aria-controls="settings" role="tab" data-toggle="tab">Settings</a></li>
            </ul>
            <div class="tab-content">
                <div role="tabpanel" class="tab-pane active" id="passwords">
                    <table class="table table-striped table-hover">
                        <thead>
                            <tr>
                                <th>Site</th>
                                <th>Password</th>
                                <th>&nbsp;</th>
                            </tr>
                        </thead>
                        <tbody>
							{{ range .Sites }}
                            <tr>
                                <td>{{ .Host }}</td>
                                <td>
                                    <button class="btn btn-default btn-xs"><span class="glyphicon glyphicon-share" aria-hidden="true"></span></button>
                                    <form class="form form-horizontal" style="display: inline-block;" action="/api/refresh" method="post">
                                        <input name="profile" type="hidden" value="{{ $profile }}" />
                                        <input name="p" type="hidden" value="{{ $passphrase }}" />
                                        <input name="host" type="hidden" value="{{ .Host }}" />
                                        <button class="btn btn-default btn-xs"><span class="glyphicon glyphicon-refresh" aria-hidden="true"></span></button>
                                    </form>
                                    <span>{{ .Password }}</span>
                                </td>
                                <td class="text-right">
                                    <form class="form form-horizontal" action="/api/remove" method="post">
                                        <input name="profile" type="hidden" value="{{ $profile }}" />
                                        <input name="p" type="hidden" value="{{ $passphrase }}" />
                                        <input name="host" type="hidden" value="{{ .Host }}" />
                                        <button class="btn btn-default"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span></button>
                                    </form>
                                </td>
                            </tr>
							{{ end }}
                        </tbody>
                    </table>
                </div>
                <div role="tabpanel" class="tab-pane" id="add">
                    <form class="form form-horizontal" action="/api/generate" method="post">
						<input name="profile" type="hidden" value="{{ $profile }}" />
						<input name="p" type="hidden" value="{{ $passphrase }}" />
                        <div class="form-group">
                            <label for="host" class="col-xs-3 control-label">Site</label>
                            <div class="col-xs-9">
                                <input id="host" name="host" type="text" class="form-control" placeholder="gmail.com" />
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="minimumLength" class="col-xs-3 control-label">Minimum Length</label>
                            <div class="col-xs-3">
                                <select id="minimumLength" name="minimumLength" class="form-control">
                                    <option value="-1">No Limit</option>
                                    <option>1</option>
                                    <option>2</option>
                                    <option>3</option>
                                    <option>4</option>
                                    <option>5</option>
                                    <option>6</option>
                                    <option>7</option>
                                    <option>8</option>
                                    <option>9</option>
                                    <option>10</option>
                                </select>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="maximumLength" class="col-xs-3 control-label">Maximum Length</label>
                            <div class="col-xs-3">
                                <select id="maximumLength" name="maximumLength" class="form-control">
                                    <option value="-1">No Limit</option>
                                    <option>4</option>
                                    <option>5</option>
                                    <option>6</option>
                                    <option>7</option>
                                    <option>8</option>
                                    <option>9</option>
                                    <option>10</option>
                                    <option>11</option>
                                    <option>12</option>
                                    <option>13</option>
                                    <option>14</option>
                                    <option>15</option>
                                    <option>16</option>
                                    <option>17</option>
                                    <option>18</option>
                                    <option>19</option>
                                    <option>20</option>
                                </select>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="minimumDigits" class="col-xs-3 control-label">Minimum Digits</label>
                            <div class="col-xs-3">
                                <select id="minimumDigits" name="minimumLength" class="form-control">
                                    <option>0</option>
                                    <option>1</option>
                                    <option>2</option>
                                    <option>3</option>
                                    <option>4</option>
                                    <option>5</option>
                                    <option>6</option>
                                    <option>7</option>
                                    <option>8</option>
                                    <option>9</option>
                                    <option>10</option>
                                </select>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="minimumUppercase" class="col-xs-3 control-label">Minimum Uppercase</label>
                            <div class="col-xs-3">
                                <select id="minimumUppercase" name="minimumUppercase" class="form-control">
                                    <option>0</option>
                                    <option>1</option>
                                    <option>2</option>
                                    <option>3</option>
                                    <option>4</option>
                                    <option>5</option>
                                    <option>6</option>
                                    <option>7</option>
                                    <option>8</option>
                                    <option>9</option>
                                    <option>10</option>
                                </select>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="minimumSpecialCharacters" class="col-xs-3 control-label">Minimum Special Characters</label>
                            <div class="col-xs-3">
                                <select id="minimumSpecialCharacters" name="minimumSpecialCharacters" class="form-control">
                                    <option>0</option>
                                    <option>1</option>
                                    <option>2</option>
                                    <option>3</option>
                                    <option>4</option>
                                    <option>5</option>
                                    <option>6</option>
                                    <option>7</option>
                                    <option>8</option>
                                    <option>9</option>
                                    <option>10</option>
                                </select>
                            </div>
                        </div>
                        <div class="form-group">
                            <label for="specialCharacters" class="col-xs-3 control-label">Special Characters</label>
                            <div class="col-xs-9">
                                <input id="specialCharacters" name="specialCharacters" type="text" class="form-control" value=" !@#$%^&*()_+-=<>,." />
                            </div>
                        </div>
                        <div class="form-group">
                            <div class="col-xs-offset-3 col-xs-10">
                                <button type="submit" class="btn btn-default">Generate Password</button>
                            </div>
                        </div>
                    </form>
                </div>
                <div role="tabpanel" class="tab-pane" id="settings">
                    <div class="row">
                        <form class="form form-horizontal" action="/api/update" method="post">
							<input name="profile" type="hidden" value="{{ $profile }}" />
							<input name="p" type="hidden" value="{{ $passphrase }}" />
                            <div class="form-group">
                                <label for="newPassphrase" class="col-xs-3 control-label">New passphrase</label>
                                <div class="col-xs-9">
                                    <input id="newPassphrase" type="password" class="form-control" placeholder="" />
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="confirmPassphrase" class="col-xs-3 control-label">Confirm passphrase</label>
                                <div class="col-xs-9">
                                    <input id="confirmPassphrase" type="password" class="form-control" placeholder="" />
                                </div>
                            </div>
                            <div class="form-group">
                                <div class="col-xs-offset-3 col-xs-6">
                                    <button name="cmd" value="update" type="submit" class="btn btn-default">Update Passphrase</button>
                                </div>
                                <div class="col-xs-3 text-right">
                                    <button name="cmd" value="delete" type="submit" class="btn btn-danger">Delete Profile</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/js/bootstrap.min.js"></script>
</body>
</html>`
)
