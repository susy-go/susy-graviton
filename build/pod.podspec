Pod::Spec.new do |spec|
  spec.name         = 'Graviton'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/susy-go/susy-graviton'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Sophon Client'
  spec.source       = { :git => 'https://github.com/susy-go/susy-graviton.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Graviton.framework'

	spec.prepare_command = <<-CMD
    curl https://gravitonstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Graviton.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
