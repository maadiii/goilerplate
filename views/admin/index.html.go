// Code generated by hero.
// source: /home/maadi/workspace/goilerplate/views/admin/index.html
// DO NOT EDIT!
package adminview

import (
	"bytes"
	"goilerplate/domain/models"
	"goilerplate/infrastructure/application"

	"github.com/shiyanhui/hero"
)

func Index(ctx *application.Context, user *models.User, buffer *bytes.Buffer) {
	buffer.WriteString(`<!DOCTYPE html>
<html lang="fa">

<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>
    `)
	buffer.WriteString(`
  </title>
  <!-- Tell the browser to be responsive to screen width -->
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <!-- Font Awesome -->
  <link rel='stylesheet' id='fontawesome-css' href='https://use.fontawesome.com/releases/v5.0.1/css/all.css?ver=4.9.1'
    type='text/css' media='all' />
  <link rel="stylesheet" href="/static/admin/plugins/font-awesome/css/font-awesome.min.css">
  <!-- Ionicons -->
  <link rel="stylesheet" href="https://code.ionicframework.com/ionicons/2.0.1/css/ionicons.min.css">
  <!-- Theme style -->
  <link rel="stylesheet" href="/static/admin/css/adminlte.min.css">
  <!-- iCheck -->
  <link rel="stylesheet" href="/static/admin/plugins/iCheck/flat/blue.css">
  <!-- Morris chart -->
  <link rel="stylesheet" href="/static/admin/plugins/morris/morris.css">
  <!-- jvectormap -->
  <link rel="stylesheet" href="/static/admin/plugins/jvectormap/jquery-jvectormap-1.2.2.css">
  <!-- Date Picker -->
  <link rel="stylesheet" href="/static/admin/plugins/datepicker/datepicker3.css">
  <!-- Daterange picker -->
  <link rel="stylesheet" href="/static/admin/plugins/daterangepicker/daterangepicker-bs3.css">
  <!-- bootstrap wysihtml5 - text editor -->
  <link rel="stylesheet" href="/static/admin/plugins/bootstrap-wysihtml5/bootstrap3-wysihtml5.min.css">
  <!-- Google Font: Source Sans Pro -->
  <link href="/static/admin/css/googleapi.css" rel="stylesheet">
  <!-- bootstrap rtl -->
  <link rel="stylesheet" href="/static/admin/css/bootstrap-rtl.min.css">
  <!-- template rtl version -->
  <link rel="stylesheet" href="/static/admin/css/custom-style.css">
  <link rel="stylesheet" href="/static/admin/css/sweetalert.css">
</head>

<body class="hold-transition sidebar-mini">
  <div class="wrapper">
    <div id="loaderC">
      <div id="loader"></div>
    </div>

    <!-- Navbar -->
    <nav class="main-header navbar navbar-expand bg-white navbar-light border-bottom">
      <!-- Left navbar links -->
      <ul class="navbar-nav">
        <li class="nav-item">
          <a class="nav-link" data-widget="pushmenu" href="#"><i class="fa fa-bars"></i></a>
        </li>
      </ul>
    </nav>
    <!-- /.navbar -->

    <!-- Main Sidebar Container -->
    <aside class="main-sidebar sidebar-dark-primary elevation-4">

      <!-- Sidebar -->
      `)

	menuOpen := "users"
	menuActive := "users"
	submenuActive := "users"

	buffer.WriteString(`

<div class="sidebar" style="direction: ltr">
  <div style="direction: rtl">
    <!-- Sidebar user panel (optional) -->
    <div class="user-panel mt-3 pb-3 mb-3 d-flex">
      <div class="image">
        <img src="https://www.gravatar.com/avatar/52f0fbcbedee04a121cba8dad1174462?s=200&d=mm&r=g"
          class="img-circle elevation-2" alt="User Image">
      </div>
      <div class="info">
        <a href="#" class="d-block">مهدی عزیزی</a>
      </div>
    </div>

    <!-- Sidebar Menu -->
    <nav class="mt-2">
      <ul class="nav nav-pills nav-sidebar flex-column" data-widget="treeview" role="menu" data-accordion="false">
        `)
	for _, menu := range menuList {
		buffer.WriteString(`
        <li class="nav-item has-treeview `)
		if Equal(menu.EnName, menuOpen) {
			buffer.WriteString(` menu-open `)
		}
		buffer.WriteString(`">
          <a href="#" class="nav-link `)
		if Equal(menu.EnName, menuActive) {
			buffer.WriteString(` active `)
		}
		buffer.WriteString(`>">
            <i class="`)
		hero.EscapeHTML(menu.Icon, buffer)
		buffer.WriteString(`"></i>
            <p>
              `)
		hero.EscapeHTML(menu.FaName, buffer)
		buffer.WriteString(`
              <i class="right fa fa-angle-left"></i>
            </p>
          </a>
          <ul class="nav nav-treeview">
            `)
		for _, submenu := range menu.Submenues {
			buffer.WriteString(`
            <li class="nav-item">
              <a href="`)
			hero.EscapeHTML(submenu.Url, buffer)
			buffer.WriteString(`" class="nav-link `)
			if Equal(submenu.EnName, submenuActive) {
				buffer.WriteString(` active `)
			}
			buffer.WriteString(`">
                <i class="`)
			hero.EscapeHTML(submenu.Icon, buffer)
			buffer.WriteString(`"></i>
                <p>`)
			hero.EscapeHTML(submenu.FaName, buffer)
			buffer.WriteString(`</p>
              </a>
            </li>
            `)
		}
		buffer.WriteString(`
          </ul>
        </li>
        `)
	}
	buffer.WriteString(`
      </ul>
    </nav>
  </div>
</div>
`)

	buffer.WriteString(`
    </aside>

    <!-- Content Wrapper. Contains page content -->
    <div class="content-wrapper">
      <!-- Content Header (Page header) -->
      <div class="content-header">
        <div class="container-fluid">
          <div class="row mb-2">
            <div class="col-sm-6">
              <h1 class="m-0 text-dark">`)
	buffer.WriteString(`</h1>
            </div><!-- /.col -->
          </div><!-- /.row -->
        </div><!-- /.container-fluid -->
      </div>
      <!-- /.content-header -->

      <!-- Main content -->
      <section class="content">
        <div class="container-fluid">
          `)
	buffer.WriteString(`
<div class="row">
    <div class="col-lg-3 col-6">
        <!-- small box -->
        <div class="small-box bg-info">
            <div class="inner">
                <h3>۱۵۰</h3>

                <p>سفارشات جدید</p>
            </div>
            <div class="icon">
                <i class="ion ion-bag"></i>
            </div>
            <a href="#" class="small-box-footer">اطلاعات بیشتر <i class="fa fa-arrow-circle-left"></i></a>
        </div>
    </div>
    <!-- ./col -->
    <div class="col-lg-3 col-6">
        <!-- small box -->
        <div class="small-box bg-success">
            <div class="inner">
                <h3>۵۳<sup style="font-size: 20px">%</sup></h3>

                <p>افزایش امتیاز</p>
            </div>
            <div class="icon">
                <i class="ion ion-stats-bars"></i>
            </div>
            <a href="#" class="small-box-footer">اطلاعات بیشتر <i class="fa fa-arrow-circle-left"></i></a>
        </div>
    </div>
    <!-- ./col -->
    <div class="col-lg-3 col-6">
        <!-- small box -->
        <div class="small-box bg-warning">
            <div class="inner">
                <h3>۴۴</h3>

                <p>کاربران ثبت شده</p>
            </div>
            <div class="icon">
                <i class="ion ion-person-add"></i>
            </div>
            <a href="#" class="small-box-footer">اطلاعات بیشتر <i class="fa fa-arrow-circle-left"></i></a>
        </div>
    </div>
    <!-- ./col -->
    <div class="col-lg-3 col-6">
        <!-- small box -->
        <div class="small-box bg-danger">
            <div class="inner">
                <h3>۶۵</h3>

                <p>بازدید جدید</p>
            </div>
            <div class="icon">
                <i class="ion ion-pie-graph"></i>
            </div>
            <a href="#" class="small-box-footer">اطلاعات بیشتر <i class="fa fa-arrow-circle-left"></i></a>
        </div>
    </div>
    <!-- ./col -->
</div>
`)

	buffer.WriteString(`
            </div><!-- /.container-fluid -->
            </section>
            <!-- /.content -->
          </div>
          <!-- /.content-wrapper -->
        </div>

        <!-- ./wrapper -->

        <!-- jQuery -->
        <script src="/static/admin/plugins/jquery/jquery.min.js"></script>
        <!-- jQuery UI 1.11.4 -->
        <script src="/static/admin/plugins/jQueryUI/jquery-ui.min.js"></script>
        <!-- Resolve conflict in jQuery UI tooltip with Bootstrap tooltip -->
        <script>
          $.widget.bridge('uibutton', $.ui.button)
        </script>
        <!-- Bootstrap 4 -->
        <script src="/static/admin/plugins/bootstrap/js/bootstrap.bundle.min.js"></script>
        <!-- Morris.js charts -->
        <script src="/static/admin/js/plugins/raphael-min.js"></script>
        <script src="/static/admin/plugins/morris/morris.min.js"></script>
        <!-- Sparkline -->
        <script src="/static/admin/plugins/sparkline/jquery.sparkline.min.js"></script>
        <!-- jvectormap -->
        <script src="/static/admin/plugins/jvectormap/jquery-jvectormap-1.2.2.min.js"></script>
        <script src="/static/admin/plugins/jvectormap/jquery-jvectormap-world-mill-en.js"></script>
        <!-- jQuery Knob Chart -->
        <script src="/static/admin/plugins/knob/jquery.knob.js"></script>
        <!-- daterangepicker -->
        <script src="/static/admin/js/plugins/moment.min.js"></script>
        <script src="/static/admin/plugins/daterangepicker/daterangepicker.js"></script>
        <!-- datepicker -->
        <script src="/static/admin/plugins/datepicker/bootstrap-datepicker.js"></script>
        <!-- Bootstrap WYSIHTML5 -->
        <script src="/static/admin/plugins/bootstrap-wysihtml5/bootstrap3-wysihtml5.all.min.js"></script>
        <!-- Slimscroll -->
        <script src="/static/admin/plugins/slimScroll/jquery.slimscroll.min.js"></script>
        <!-- FastClick -->
        <script src="/static/admin/plugins/fastclick/fastclick.js"></script>
        <!-- AdminLTE App -->
        <script src="/static/admin/js/adminlte.js"></script>
        <!-- AdminLTE dashboard demo (This is only for demo purposes) -->
        <script src="/static/admin/js/pages/dashboard.js"></script>
        <!-- AdminLTE for demo purposes -->
        <script src="/static/admin/js/demo.js"></script>
        <script src="/static/admin/js/sweetalert2.js"></script>
        <script src="/static/admin/js/custom.js"></script>
        `)
	buffer.WriteString(`

</body>

</html>
`)

}
